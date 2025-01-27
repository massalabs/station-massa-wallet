*** Settings ***
Documentation       This is a test suite for Massa Wallet /sign endpoints.

Library             RequestsLibrary
Library             Collections
Resource            ../variables.resource


*** Test Cases ***
POST Sign CallSC
    ${data}=    Create Dictionary
    ...    description=this is the description
    ...    operation=lQbAxAcEwIQ9uWAAAC5IN91TH6nQKZLSUhsico2f9dG9KI1+e5zfu81A7Ci4D2V4YW1wbGVGdW5jdGlvbhFleGFtcGxlUGFyYW1ldGVycw==
    ...    batch=${false}
    ...    chainId=${CHAIN_ID}
    ${response}=    POST    ${API_URL}/sign    json=${data}    expected_status=any
    Log To Console    json response: ${response.json()}    # Print the response content to the test log for debugging
    Should Be Equal As Integers    ${response.status_code}    ${STATUS_OK}    # Assert the status code is 200 OK

