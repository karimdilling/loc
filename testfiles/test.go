// This is a file to run the tests for the countLines() function on
package main

import "fmt"

/* Added a comment here */
func main() {
	fmt.Println("Hello world! /*") // also ignore "/*" inside a string when parsing for comments
	fmt.Println("*/ Test")
	fmt.Println("/*")
	fmt.Println("'/*'")
	abc()
}

/* A multiline comment
is done here */

func abc() int {
	return 100 * 2
} /* A multline comment can also start with
some code at the beginning
*/

/* Or involve some asterisks inside
* it ** like * this
right here **/
