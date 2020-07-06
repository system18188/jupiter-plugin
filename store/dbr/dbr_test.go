package dbr

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/system18188/jupiter-plugin/store/dbr/dialect"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mailru/go-clickhouse"
	_ "github.com/mattn/go-sqlite3"
)

// Ensure that tx and session are session runner
var (
	_ SessionRunner = (*Tx)(nil)
	_ SessionRunner = (*Session)(nil)
)

var (
	currID int64 = 256
)

// nextID returns next pseudo unique id
func nextID() int64 {
	currID++
	return currID
}

const (
	mysqlDSN      = "root@unix(/tmp/mysql.sock)/dbr_test?charset=utf8"
	postgresDSN   = "postgres://postgres@localhost:5432/dbr_test?sslmode=disable"
	sqlite3DSN    = ":memory:"
	clickhouseDSN = "http://localhost:8123/dbr_test"
)

func createSession(driver, dsn string) *Session {
	var testDSN string
	switch driver {
	case "mysql":
		testDSN = os.Getenv("DBR_TEST_MYSQL_DSN")
	case "postgres":
		testDSN = os.Getenv("DBR_TEST_POSTGRES_DSN")
	case "sqlite3":
		testDSN = os.Getenv("DBR_TEST_SQLITE3_DSN")
	case "clickhouse":
		testDSN = os.Getenv("DBR_TEST_CLICKHOUSE_DSN")
	}
	if testDSN != "" {
		dsn = testDSN
	}
	conn, err := Open(driver, dsn, nil)
	if err != nil {
		log.Fatal(err)
	}
	sess := conn.NewSession(nil)
	reset(sess)
	return sess
}

var (
	mysqlSession          = createSession("mysql", mysqlDSN)
	postgresSession       = createSession("postgres", postgresDSN)
	postgresBinarySession = createSession("postgres", postgresDSN+"&binary_parameters=yes")
	sqlite3Session        = createSession("sqlite3", sqlite3DSN)
	clickhouseSession     = createSession("clickhouse", clickhouseDSN)

	// all test sessions should be here
	testSession = []*Session{mysqlSession, postgresSession, sqlite3Session, clickhouseSession}
)

type person struct {
	ID    int64
	Name  string
	Email string
}

type nullTypedRecord struct {
	ID         int64
	StringVal  NullString
	Int64Val   NullInt64
	Float64Val NullFloat64
	TimeVal    NullTime
	BoolVal    NullBool
}

func reset(sess *Session) {
	var stmts []string
	switch sess.Dialect {
	case dialect.MySQL:
		stmts = []string{
			`DROP TABLE IF EXISTS dbr_people`,
			`CREATE TABLE dbr_people(id SERIAL PRIMARY KEY, name varchar(255) NOT NULL, email varchar(255))`,
			`DROP TABLE IF EXISTS null_types`,
			`CREATE TABLE null_types(
				id SERIAL PRIMARY KEY,
				string_val varchar(255) NULL,
				int64_val integer NULL,
				float64_val float NULL,
				time_val timestamp NULL,
				bool_val bool NULL
			)`,
			`DROP TABLE IF EXISTS dbr_keys`,
			`CREATE TABLE dbr_keys (key_value varchar(255) PRIMARY KEY, val_value varchar(255))`,
		}
	case dialect.PostgreSQL:
		stmts = []string{
			"DROP TABLE IF EXISTS dbr_people",
			"CREATE TABLE dbr_people(id SERIAL PRIMARY KEY, name varchar(255) NOT NULL, email varchar(255))",
			`DROP TABLE IF EXISTS null_types`,
			`CREATE TABLE null_types(
				id SERIAL PRIMARY KEY,
				string_val varchar(255) NULL,
				int64_val integer NULL,
				float64_val float NULL,
				time_val timestamp NULL,
				bool_val bool NULL
			)`,
			`DROP TABLE IF EXISTS dbr_keys`,
			`CREATE TABLE dbr_keys (key_value varchar(255) PRIMARY KEY, val_value varchar(255))`,
		}
	case dialect.SQLite3:
		stmts = []string{
			"DROP TABLE IF EXISTS dbr_people",
			"CREATE TABLE dbr_people(id INTEGER PRIMARY KEY, name varchar(255) NOT NULL, email varchar(255))",
			`DROP TABLE IF EXISTS null_types`,
			`CREATE TABLE null_types(
				id INTEGER PRIMARY KEY,
				string_val varchar(255) NULL,
				int64_val integer NULL,
				float64_val float NULL,
				time_val timestamp NULL,
				bool_val bool NULL
			)`,
			`DROP TABLE IF EXISTS dbr_keys`,
			`CREATE TABLE dbr_keys (key_value varchar(255) PRIMARY KEY, val_value varchar(255))`,
		}
	case dialect.ClickHouse:
		stmts = []string{
			"DROP TABLE IF EXISTS dbr_people",
			"CREATE TABLE dbr_people(id Int32, name String, email String) Engine=Memory",
			`DROP TABLE IF EXISTS dbr_keys`,
			`CREATE TABLE dbr_keys (key_value String, val_value String) Engine=Memory`,
		}
	}
	for _, v := range stmts {
		_, err := sess.Exec(v)
		if err != nil {
			log.Fatalf("Failed to execute statement: %s, Got error: %s", v, err)
		}
	}
}

func BenchmarkByteaNoBinaryEncode(b *testing.B) {
	benchmarkBytea(b, postgresSession)
}

func BenchmarkByteaBinaryEncode(b *testing.B) {
	benchmarkBytea(b, postgresBinarySession)
}

func benchmarkBytea(b *testing.B, sess *Session) {
	data := bytes.Repeat([]byte("0123456789"), 1000)
	for _, v := range []string{
		`DROP TABLE IF EXISTS bytea_table`,
		`CREATE TABLE bytea_table (
			val bytea
		)`,
	} {
		_, err := sess.Exec(v)
		assert.NoError(b, err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := sess.InsertInto("bytea_table").Pair("val", data).Exec()
		assert.NoError(b, err)
	}
}

func TestBasicCRUD(t *testing.T) {
	for _, sess := range testSession {
		jonathan := person{
			Name:  "jonathan",
			Email: "jonathan@uservoice.com",
		}
		insertColumns := []string{"name", "email"}
		if sess.Dialect == dialect.PostgreSQL || sess.Dialect == dialect.ClickHouse {
			jonathan.ID = nextID()
			insertColumns = []string{"id", "name", "email"}
		}
		// insert
		result, err := sess.InsertInto("dbr_people").Columns(insertColumns...).Record(&jonathan).Exec()
		assert.NoError(t, err)

		rowsAffected, err := result.RowsAffected()
		if err == nil {
			assert.EqualValues(t, 1, rowsAffected)
		}

		assert.True(t, jonathan.ID > 0)
		// select
		var people []person
		count, err := sess.Select("*").From("dbr_people").Where(Eq("id", jonathan.ID)).LoadStructs(&people)
		assert.NoError(t, err)
		if assert.Equal(t, 1, count) {
			assert.Equal(t, jonathan.ID, people[0].ID)
			assert.Equal(t, jonathan.Name, people[0].Name)
			assert.Equal(t, jonathan.Email, people[0].Email)
		}

		// select id
		ids, err := sess.Select("id").From("dbr_people").ReturnInt64s()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(ids))

		// select id limit
		ids, err = sess.Select("id").From("dbr_people").Limit(1).ReturnInt64s()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(ids))

		if sess.Dialect == dialect.ClickHouse {
			// clickhouse does not support update/delete
			continue
		}
		// update
		result, err = sess.Update("dbr_people").Where(Eq("id", jonathan.ID)).Set("name", "jonathan1").Exec()
		assert.NoError(t, err)

		rowsAffected, err = result.RowsAffected()
		assert.NoError(t, err)
		assert.EqualValues(t, 1, rowsAffected)

		var n NullInt64
		sess.Select("count(*)").From("dbr_people").Where("name = ?", "jonathan1").LoadValue(&n)
		assert.EqualValues(t, 1, n.Int64)

		// delete
		result, err = sess.DeleteFrom("dbr_people").Where(Eq("id", jonathan.ID)).Exec()
		assert.NoError(t, err)

		rowsAffected, err = result.RowsAffected()
		assert.NoError(t, err)
		assert.EqualValues(t, 1, rowsAffected)

		// select id
		ids, err = sess.Select("id").From("dbr_people").ReturnInt64s()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(ids))
	}
}

func TestOnConflict(t *testing.T) {
	for _, sess := range testSession {
		if sess.Dialect == dialect.SQLite3 || sess.Dialect == dialect.ClickHouse {
			continue
		}
		for i := 0; i < 2; i++ {
			b := sess.InsertInto("dbr_keys").Columns("key_value", "val_value").Values("key", "value")
			b.OnConflict("dbr_keys_pkey").Action("val_value", Expr("CONCAT(?, 2)", Proposed("val_value")))
			_, err := b.Exec()
			assert.NoError(t, err)
		}
		var value string
		_, err := sess.SelectBySql("SELECT val_value FROM dbr_keys WHERE key_value=?", "key").Load(&value)
		assert.NoError(t, err)
		assert.Equal(t, "value2", value)
	}
}

func TestForkSession(t *testing.T) {
	sess := testSession[0]
	sess2 := sess.NewSession(nil)
	assert.True(t, sess.ctx == sess2.ctx)
	assert.True(t, sess.EventReceiver == sess2.EventReceiver)
	recv := new(NullEventReceiver)
	sess3 := sess.NewSession(recv)
	assert.True(t, sess3.EventReceiver != sess.EventReceiver)
}
