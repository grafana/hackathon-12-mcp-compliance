# Manual Testing Guide

This document contains example questions to test the FedRAMP Compliance MCP Server functionality.
Use these questions to verify that the server is working correctly and providing accurate information.

## Test Questions

### General Program Questions

1. "What security controls are required for FedRAMP High?"
2. "What security controls are required for FedRAMP Moderate?"
3. "What are the differences between FedRAMP High and Moderate for access control requirements?"
4. "Which program should I use for my cloud service?"

### Control-Specific Questions

1. "What are the requirements for access control in FedRAMP High?"
2. "How does FedRAMP Moderate handle incident response?"
3. "What is control AC-2 in FedRAMP High?"
4. "What parameters are required for AC-2 in FedRAMP Moderate?"
5. "Show me all the controls in the AC family for FedRAMP High"
6. "What are the audit requirements in FedRAMP Moderate?"

### Implementation Questions

1. "How do I implement AC-2 for FedRAMP High?"
2. "What evidence do I need to collect for AC-1 in FedRAMP Moderate?"
3. "What are the password requirements for FedRAMP High?"
4. "How often do I need to review access in FedRAMP Moderate?"
5. "What documentation is required for incident response in FedRAMP High?"

### Comparison Questions

1. "How do the access control requirements differ between FedRAMP High and Moderate?"
2. "Which program has stricter password requirements?"

## Expected Results

For each question, verify that:
1. The server returns a response in a reasonable time
2. The response is well-formatted JSON
3. The information is accurate according to FedRAMP documentation
4. All required fields are present in the response
5. Error cases are handled gracefully

## Common Test Scenarios

### Control Family Tests
1. Request all controls in the AC (Access Control) family
2. Request all controls in the AU (Audit) family
3. Request a non-existent control family

### Control Parameter Tests
1. Request parameters for AC-2
2. Request parameters for a control without parameters
3. Request parameters for a non-existent control

### Search Tests
1. Search for controls containing "password"
2. Search for controls containing "audit"
3. Search with an empty query
4. Search with special characters

### Program Comparison Tests
1. Compare AC-2 between High and Moderate
2. Compare AU-1 between High and Moderate
3. Compare a control that exists in one program but not the other 