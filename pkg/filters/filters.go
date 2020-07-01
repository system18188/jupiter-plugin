package filters

import (
	"regexp"
)

var (
	CompileHtml   = regexp.MustCompile(`\<[\S\s]+?\>`)
	CompileStyle  = regexp.MustCompile(`\<style[\S\s]+?\</style\>`)
	CompileScript = regexp.MustCompile(`\<script[\S\s]+?\</script\>`)
	CompileS2     = regexp.MustCompile(`\s{2,}`)
	CompileUrl    = regexp.MustCompile(`((ht|f)tps?):\/\/([\w\-]+(\.[\w\-]+)*\/)*[\w\-]+(\.[\w\-]+)*\/?(\?([\w\-\.,@?^=%&:\/~\+#]*)+)?`)
)
