# Package check

Package check contains test helper functions.

Currently, it only contains two functions: 
* ErrorString
* Error

# ErrorString

ErrorString compares the message of an error (err) 
with a provided string (want).

It is a convenience method for simplifying error tests
where the comparison of the error messages is sufficient.
Ensure that the wanted error message is long enough to build
a valid test case.

Do not use it when custom and standard error types have to be 
compared exactly (for example: TestCasterValidate in 
[caster\_test.go](../alchemist/caster_test.go)).

It is a test helper function and calls t.Helper() at the beginning.
The effect is that the location of the error messages in the go test
report is the lines where ErrorString has been called, not the lines 
within ErrorString where the message has been generated.

It behaves like this:
* If want is empty and err is nil, no test error is raised.
* If want is not empty err is not nil and it is a subset of the error message, no test error is raised (the error message is displayed when go test runs in verbose mode)
* Otherwise: a Fatal test error is raised.

Additionally:
* If an error is expected (that is: want is not empty), the the rest of the test is skipped.
* If go test runs in verbose mode, the skipping is reported.

Skip means that the statements in the unit test after the call to 
ErrorString are not executed, because usually you do not expect a valid result
in an error case
(for example: TestListPages in [book\_test.go](../alchemist/book_test.go)).

In some cases where values should be tested even in error cases, just place
the call to ErrorString at the end of the test function
(for example: TestSpellCast in [spell\_test.go](../alchemist/spell_test.go)).


# Testing ErrorString

To be able to test the ErrorString function, it does not take a real
testing.T object, but an interface that is defined in the check package.
It only provides the subset of methods from testing.T that are used
in ErrorString. 

This is a typical situation where it comes very handy that go does
not require that a type declares the interfaces it implements.
testing.T does not know it implements the interface "errStringTester",
but it does because its methods match all the methods of the interface.

So the function can be used with testing.T in the other
unit tests and with a testerSpy when it is tested itself in
[errorstring\_test.go](errorstring_test.go).

# Error

Error does pretty much the same as ErrorString, but the expectation
is provided as an error value. It is also testable.

