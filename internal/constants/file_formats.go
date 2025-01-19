package constants

const GOLANG_FILE string = `
/*
%s

(NOTE: The setup for the code may not be exact)
Submit your solution at: %s
*/
package main

func %s(/*ENTER PARAMS*/) /*RETURN TYPE*/ {
        
}
`

const JAVA_FILE string = `
/*
%s

(NOTE: The setup for the code may not be exact)
Submit your solution at: %s
*/
class Solution {
    public /*RETURN TYPE*/ %s(/*ENTER PARAMS*/) {
        
    }
}
`
const PYTHON_FILE string = `
"""
%s

(NOTE: The setup for the code may not be exact)
Submit your solution at: %s
"""
class Solution:
	def %s(self, """add params here"""):
		raise NotImplemented()
`