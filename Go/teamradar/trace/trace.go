/*
Copyright 2016 IslandJohn and the TeamRadar Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package trace

import (
	"fmt"
	"log"
	"runtime"
)

// http://stackoverflow.com/questions/25927660/golang-get-current-scope-of-function-name
func Where(n int) (string, string, int) {
	pc := make([]uintptr, n+1) // at least 1 entry needed
	runtime.Callers(n+2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])

	return file, f.Name(), line
}

func Here() (string, string, int) {
	return Where(1)
}

func Log(v ...interface{}) {
	_, f, line := Where(1)
	log.Printf("%s:%d %s", f, line, fmt.Sprint(v...))
}

func Logf(format string, v ...interface{}) {
	_, f, line := Where(1)
	log.Printf("%s:%d %s", f, line, fmt.Sprintf(format, v...))
}
