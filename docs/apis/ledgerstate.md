# Ledgerstate API Methods 
## HTTP APIs:

* [/ledgerstate/addresses/:address](#ledgerstateaddressesaddress)
* [/ledgerstate/addresses/:address/unspentOutputs](#ledgerstateaddressesaddressunspentoutputs)
* [/ledgerstate/branches/:branchID](#ledgerstatebranchesbranchid)
* [/ledgerstate/branches/:branchID/children](#ledgerstatebranchesbranchidchildren)
* [/ledgerstate/branches/:branchID/conflicts](#ledgerstatebranchesbranchidconflicts)
* [/ledgerstate/outputs/:outputID](#ledgerstateoutputsoutputid)
* [/ledgerstate/outputs/:outputID/consumers](#ledgerstateoutputsoutputidconsumers)
* [/ledgerstate/outputs/:outputID/metadata](#ledgerstateoutputsoutputidmetadata)
* [/ledgerstate/transactions/:transactionID](#ledgerstatetransactionstransactionid)
* [/ledgerstate/transactions/:transactionID/metadata](#ledgerstatetransactionstransactionidmetadata)
* [/ledgerstate/transactions/:transactionID/inclusionState](#ledgerstatetransactionstransactionidinclusionstate)
* [/ledgerstate/transactions/:transactionID/consensus](#ledgerstatetransactionstransactionidconsensus)
* [/ledgerstate/transactions/:transactionID/attachments](#ledgerstatetransactionstransactionidattachments)
* [/ledgerstate/transactions](#ledgerstatetransactions)
* [/ledgerstate/addresses/unspentOutputs](#ledgerstateaddressesunspentoutputs)


## Client lib APIs:
* [GetAddressOutputs()](#client-lib---getaddressoutputs)
* [GetAddressUnspentOutputs()](#client-lib---getaddressunspentoutputs)
* [GetBranch()](#client-lib---getbranch)
* [GetBranchChildren()](#client-lib---getbranchchildren)
* [GetBranchConflicts()](#client-lib---getbranchconflicts)
* [GetOutput()](#client-lib---getoutput)
* [GetOutputConsumers()](#client-lib---getoutputconsumers)
* [GetOutputMetadata()](#client-lib---getoutputmetadata)
* [GetTransaction()](#client-lib---gettransaction)
* [GetTransactionMetadata()](#client-lib---gettransactionmetadata)
* [GetTransactionInclusionState()](#client-lib---gettransactioninclusionstate)
* [GetTransactionConsensusMetadata()](#client-lib---gettransactionconsensusmetadata)
* [GetTransactionAttachments()](#client-lib---gettransactionattachments)
* [PostTransaction()](#client-lib---posttransaction)
* [PostAddressUnspentOutputs()](#client-lib---postaddressunspentoutputs)

## `/ledgerstate/addresses/:address`
Get address details for a given base58 encoded address ID, such as output types and balances. For the client library API call balances will not be directly available as values because they are stored as a raw message. Balance can be read after retrieving `ledgerstate.Output` instance, as presented in the examples.

### Parameters
| **Parameter**            | `address`      |
|--------------------------|----------------|
| **Required or Optional** | required       |
| **Description**          | The address encoded in base58. |
| **Type**                 | string         |

### Examples

#### cURL

```shell
curl http://localhost:8080/ledgerstate/addresses/:address \
-X GET \
-H 'Content-Type: application/json'
```

where `:address` is the base58 encoded address, e.g. 6PQqFcwarCVbEMxWFeAqj7YswK842dMtf84qGyKqVH7s1kK.

#### Client lib - `GetAddressOutputs()`
```Go
resp, err := goshimAPI.GetAddressOutputs("6PQqFcwarCVbEMxWFeAqj7YswK842dMtf84qGyKqVH7s1kK")
if err != nil {
    // return error
}
fmt.Println("output address: ", resp.Address)

for _, output := range resp.Outputs {
    fmt.Println("outputID: ", output.OutputID)
    fmt.Println("output type: ", output.Type)
    // get output instance
    out, err = output.ToLedgerstateOutput()
}
```

### Response examples
```json
{
    "address": {
        "type": "AddressTypeED25519",
        "base58": "18LhfKUkWt4M9YR6Q3au4LT8wWCERwzHaqn153K78Eixp"
    },
    "outputs": [
        {
            "outputID": {
                "base58": "gdFXAjwsm5kDeGdcZsJAShJLeunZmaKEMmfHSdoX34ZeSs",
                "transactionID": "32yHjeZpghKNkybd2iHjXj7NsUdR63StbJcBioPGAut3",
                "outputIndex": 0
            },
            "type": "SigLockedColoredOutputType",
            "output": {
                "balances": {
                    "11111111111111111111111111111111": 1000000
                },
                "address": "18LhfKUkWt4M9YR6Q3au4LT8wWCERwzHaqn153K78Eixp"
            }
        }
    ]
}
```

### Results
|Return field | Type | Description|
|:-----|:------|:------|
| `address`  | Address | The address corresponding to provided outputID.   |
| `outputs`   | Output | List of transactions' outputs.     |

#### Type `Address`

|Field | Type | Description|
|:-----|:------|:------|
| `type`  | string | The type of an address.   |
| `base58`| string   | The address encoded with base58.          |

#### Type `Output`

|Field | Type | Description|
|:-----|:------|:------|
| `outputID`  | OutputID | The identifier of an output.   |
| `outputType`   | string | The type of the output.     |
| `output` | string | An output raw message containing balances and corresponding addresses. |

#### Type `OutputID`

|Field | Type | Description|
|:-----|:------|:------|
| `base58`  | string | The output identifier encoded with base58.    |
| `transactionID`   | string | The transaction identifier encoded with base58.     |
| `outputIndex`   | int | The index of an output.     |

<br />

## `/ledgerstate/addresses/:address/unspentOutputs`
Gets list of all unspent outputs for the address based on a given base58 encoded address ID.

### Parameters

| **Parameter**            | `address`      |
|--------------------------|----------------|
| **Required or Optional** | required       |
| **Description**          | The address encoded in base58. |
| **Type**                 | string         |
### Examples

#### cURL

```shell
curl http://localhost:8080/ledgerstate/addresses/:address/unspentOutputs \
-X GET \
-H 'Content-Type: application/json'
```

where `:address` is the base58 encoded address, e.g. 6PQqFcwarCVbEMxWFeAqj7YswK842dMtf84qGyKqVH7s1kK.

#### Client lib - `GetAddressUnspentOutputs()`

```Go
address := "6PQqFcwarCVbEMxWFeAqj7YswK842dMtf84qGyKqVH7s1kK"
resp, err := goshimAPI.GetAddressUnspentOutputs(address)
if err != nil {
    // return error
}
fmt.Println("output address: ", resp.Address)

for _, output := range resp.Outputs {
    fmt.Println("outputID: ", output.OutputID)
    fmt.Println("output type: ", output.Type)
    // get output instance
    out, err = output.ToLedgerstateOutput()
}
```
### Response examples
```json
{
    "address": {
        "type": "AddressTypeED25519",
        "base58": "18LhfKUkWt4M9YR6Q3au4LT8wWCERwzHaqn153K78Eixp"
    },
    "outputs": [
        {
            "outputID": {
                "base58": "gdFXAjwsm5kDeGdcZsJAShJLeunZmaKEMmfHSdoX34ZeSs",
                "transactionID": "32yHjeZpghKNkybd2iHjXj7NsUdR63StbJcBioPGAut3",
                "outputIndex": 0
            },
            "type": "SigLockedColoredOutputType",
            "output": {
                "balances": {
                    "11111111111111111111111111111111": 1000000
                },
                "address": "18LhfKUkWt4M9YR6Q3au4LT8wWCERwzHaqn153K78Eixp"
            }
        }
    ]
}
```

### Results
|Return field | Type | Description|
|:-----|:------|:------|
| `address`  | Address | The address corresponding to provided unspent outputID.   |
| `outputs`   | Output | List of transactions' unspent outputs.     |

#### Type `Address`

|Field | Type | Description|
|:-----|:------|:------|
| `type`  | string | The type of an address.   |
| `base58`| string   | The address encoded with base58.          |

#### Type `Output`

|Field | Type | Description|
|:-----|:------|:------|
| `outputID`  | OutputID | The identifier of an output.   |
| `outputType`   | string | The type of the output.    |
| `output` | string | An output raw message containing balances and corresponding addresses |

#### Type `OutputID`

|Field | Type | Description|
|:-----|:------|:------|
| `base58`  | string | The output identifier encoded with base58.    |
| `transactionID`   | string | The transaction identifier encoded with base58.     |
| `outputIndex`   | int | The index of an output.     |

<br />

## `/ledgerstate/branches/:branchID`
Gets a branch details for a given base58 encoded branch ID.

### Parameters

| **Parameter**            | `branchID`      |
|--------------------------|----------------|
| **Required or Optional** | required       |
| **Description**          | The branch ID encoded in base58. |
| **Type**                 | string         |

### Examples

#### cURL

```shell
curl http://localhost:8080/ledgerstate/branches/:branchID \
-X GET \
-H 'Content-Type: application/json'
```

where `:branchID` is the ID of the branch, e.g. 2e2EU6fhxRhrXVnYQ6US4zmUkE5YJip25ecafn8gZeoZ.

#### Client lib - `GetBranch()`
```Go
resp, err := goshimAPI.GetBranch("2e2EU6fhxRhrXVnYQ6US4zmUkE5YJip25ecafn8gZeoZ")
if err != nil {
    // return error
}
fmt.Println("branch ID: ", resp.ID)
fmt.Println("branch type: ", resp.Type)
fmt.Println("branch inclusion state: ", resp.InclusionState)
fmt.Println("branch parents IDs: ", resp.Parents)
fmt.Println("branch conflicts IDs: ", resp.ConflictIDs)
fmt.Printf("liked: %v, finalized: %v, monotonically liked: %v", resp.Liked, resp.Finalized, resp.MonotonicallyLiked)
```
### Response examples
```json
{
    "id": "5v6iyxKUSSF73yoZa6YngNN5tqoX8hJQWKGXrgcz3XTg",
    "type": "ConflictBranchType",
    "parents": [
        "4uQeVj5tqViQh7yWWGStvkEG1Zmhx6uasJtWCJziofM"
    ],
    "conflictIDs": [
        "3LrHecDf8kvDGZKTAYaKmvdsqXA18YBc8A9UePu7pCxw5ks"
    ],
    "liked": false,
    "monotonicallyLiked": false,
    "finalized": false,
    "inclusionState": "InclusionState(Pending)"
}
```

### Results
|Return field | Type | Description|
|:-----|:------|:------|
| `id`  | string | The branch identifier encoded with base58.   |
| `type`        | string | The type of the branch.  |
| `parents`     | []string | The list of parent branches IDs.   |
| `conflictIDs` | []string | The list of conflicts identifiers.  |
| `liked`       | bool | The boolean indicator if branch is liked.   |
| `monotonicallyLiked`   | bool | The boolean indicator if branch is monotonically liked.   |
| `finalized`   | bool | The boolean indicator if branch is finalized.   |
| `inclusionState`   | string | Inclusion state of a branch.   |

<br />

## `/ledgerstate/branches/:branchID/children`
Gets a list of all child branches for a branch with given base58 encoded branch ID.

### Parameters

| **Parameter**            | `branchID`      |
|--------------------------|----------------|
| **Required or Optional** | required       |
| **Description**          | The branch ID encoded in base58. |
| **Type**                 | string         |


### Examples

#### cURL

```shell
curl http://localhost:8080/ledgerstate/branches/:branchID/children \
-X GET \
-H 'Content-Type: application/json'
```

where `:branchID` is the ID of the branch, e.g. 2e2EU6fhxRhrXVnYQ6US4zmUkE5YJip25ecafn8gZeoZ.

#### Client lib - `GetBranchChildren()`
```Go
resp, err := goshimAPI.GetBranchChildren("2e2EU6fhxRhrXVnYQ6US4zmUkE5YJip25ecafn8gZeoZ")
if err != nil {
    //return error
}
fmt.Printf("All children branches for branch %s:\n", resp.BranchID)
for _, branch := range resp.ChildBranches {
    fmt.Println("branchID: ", branch.BranchID)
    fmt.Printf("type: %s\n", branch.BranchID)
}
```

### Response examples
```json
{
    "branchID": "HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV",
    "childBranches": [
        {
            "branchID": "4SdXm5NXEcVogiJNEKkecqd5rZzRYeGYBj8oBNsdX91W",
            "type": "AggregatedBranchType"
        }
    ]
}
```

### Results

|Return field | Type | Description|
|:-----|:------|:------|
| `branchID`  | string | The branch identifier encoded with base58.   |
| `childBranches`        | []ChildBranch | The child branches data.  |


#### Type `ChildBranch`

|Field | Type | Description|
|:-----|:------|:------|
| `branchID`  | string | The branch identifier encoded with base58.   |
| `type`        | string | The type of the branch.  |

<br />

## `/ledgerstate/branches/:branchID/conflicts`
Get all conflicts for a given branch ID, their outputs and conflicting branches.

### Parameters

| **Parameter**            | `branchID`      |
|--------------------------|----------------|
| **Required or Optional** | required       |
| **Description**          | The conflicting branch ID encoded in base58. |
| **Type**                 | string         |


### Examples

#### cURL

```shell
curl http://localhost:8080/ledgerstate/branches/:branchID/conflicts \
-X GET \
-H 'Content-Type: application/json'
```

where `:branchID` is the ID of the branch, e.g. 2e2EU6fhxRhrXVnYQ6US4zmUkE5YJip25ecafn8gZeoZ.

#### Client lib - `GetBranchConflicts()`
```Go
resp, err := goshimAPI.GetBranchConflicts("2e2EU6fhxRhrXVnYQ6US4zmUkE5YJip25ecafn8gZeoZ")
if err != nil {
    // return error
}
fmt.Printf("All conflicts for branch %s:\n", resp.BranchID)
// iterate over all conflicts
for _, branch := range resp.Conflicts {
    fmt.Println("output ID: ", branch.OutputID.Base58)
    fmt.Println("conflicting transaction ID: ", branch.OutputID.TransactionID)
    fmt.Printf("related branches: %v\n", branch.BranchIDs)
}
```
### Response examples
```json
{
    "branchID": "HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV",
    "conflicts": [
        {
            "outputID": {
                "base58": "41GvDSQnd12e4nWnd2WzmdLmffruXqsE46jgeUbnB8s1QnK",
                "transactionID": "9wr21zza46Y5QonKEHNQ6x8puA7Rbq5LAbsQZJCK1g1g",
                "outputIndex": 0
            },
            "branchIDs": [
                "b8QRhHerfg14cYQ4VFD7Fyh1HYTCbjt9aK1XJmdoXwq",
                "HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV"
            ]
        }
    ]
}
```

### Results

|Return field | Type | Description|
|:-----|:------|:------|
| `branchID`  | string | The branch identifier encoded with base58.   |
| `conflicts` | []Conflict | The conflict data.  |

#### Type `Conflict`
|Field | Type | Description|
|:-----|:------|:------|
| `outputID`  | OutputID | The branch identifier encoded with base58.   |
| `branchIDs` | []string | The identifiers of all related branches encoded in base58.  |

#### Type `OutputID`

|Field | Type | Description|
|:-----|:------|:------|
| `base58`  | string | The output identifier encoded with base58.    |
| `transactionID`   | string | The transaction identifier encoded with base58.     |
| `outputIndex`   | int | The index of an output.     |

<br />

## `/ledgerstate/outputs/:outputID`
Get an output details for a given base58 encoded output ID, such as output types, addresses, and their corresponding balances.
For the client library API call balances will not be directly available as values because they are stored as a raw message. 

### Parameters

| **Parameter**            | `outputID`      |
|--------------------------|----------------|
| **Required or Optional** | required       |
| **Description**          | The output ID encoded in base58. |
| **Type**                 | string         |


### Examples

#### cURL

```shell
curl http://localhost:8080/ledgerstate/outputs/:outputID \
-X GET \
-H 'Content-Type: application/json'
```

where `:outputID` is the ID of the output, e.g. 41GvDSQnd12e4nWnd2WzmdLmffruXqsE46jgeUbnB8s1QnK.

#### Client lib - `GetOutput()`
```Go
resp, err := goshimAPI.GetOutput("41GvDSQnd12e4nWnd2WzmdLmffruXqsE46jgeUbnB8s1QnK")
if err != nil {
    // return error
}
fmt.Println("outputID: ", resp.OutputID.Base58)
fmt.Println("output type: ", resp.Type)
fmt.Println("transactionID: ", resp.OutputID.TransactionID)
```
### Response examples
```json
{
    "outputID": {
        "base58": "41GvDSQnd12e4nWnd2WzmdLmffruXqsE46jgeUbnB8s1QnK",
        "transactionID": "9wr21zza46Y5QonKEHNQ6x8puA7Rbq5LAbsQZJCK1g1g",
        "outputIndex": 0
    },
    "type": "SigLockedColoredOutputType",
    "output": {
        "balances": {
            "11111111111111111111111111111111": 1000000
        },
        "address": "1F95a2yceDicNLvqod6P3GLFZDAFdwizcTTYow4Y1G3tt"
    }
}
```

### Results
|Return field | Type | Description|
|:-----|:------|:------|
| `outputID`  | OutputID | The identifier of an output.   |
| `outputType`   | string | The type of the output.     |
| `output` | string | An output raw message containing balances and corresponding addresses |

#### Type `OutputID`

|Field | Type | Description|
|:-----|:------|:------|
| `base58`  | string | The output identifier encoded with base58.    |
| `transactionID`   | string | The transaction identifier encoded with base58.     |
| `outputIndex`   | int | The index of an output.     |

<br />

## `/ledgerstate/outputs/:outputID/consumers`
Get a list of consumers based on a provided base58 encoded output ID. Transactions that contains the output and information about its validity.

### Parameters

| **Parameter**            | `outputID`      |
|--------------------------|----------------|
| **Required or Optional** | required       |
| **Description**          | The output ID encoded in base58. |
| **Type**                 | string         |

### Examples

#### cURL

```shell
curl http://localhost:8080/ledgerstate/outputs/:outputID/consumers \
-X GET \
-H 'Content-Type: application/json'
```

where `:outputID` is the ID of the output, e.g. 41GvDSQnd12e4nWnd2WzmdLmffruXqsE46jgeUbnB8s1QnK.

#### Client lib - `GetOutputConsumers()`
```Go
resp, err := goshimAPI.GetOutputConsumers("41GvDSQnd12e4nWnd2WzmdLmffruXqsE46jgeUbnB8s1QnK")
if err != nil {
    // return error
}
fmt.Println("outputID: ", resp.OutputID.Base58)
// iterate over output consumers
for _, consumer := range resp.Consumers {
    fmt.Println("transactionID: ", consumer.TransactionID)
    fmt.Println("valid: ", consumer.Valid)
}
```
### Response examples
```json
{
    "outputID": {
        "base58": "41GvDSQnd12e4nWnd2WzmdLmffruXqsE46jgeUbnB8s1QnK",
        "transactionID": "9wr21zza46Y5QonKEHNQ6x8puA7Rbq5LAbsQZJCK1g1g",
        "outputIndex": 0
    },
    "consumers": [
        {
            "transactionID": "b8QRhHerfg14cYQ4VFD7Fyh1HYTCbjt9aK1XJmdoXwq",
            "valid": "true"
        },
        {
            "transactionID": "HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV",
            "valid": "true"
        }
    ]
}
```

### Results
|Return field | Type | Description|
|:-----|:------|:------|
| `outputID`  | OutputID | The output identifier encoded with base58.   |
| `consumers`        | []Consumer | Consumers of the requested output. |


#### Type `OutputID`

|Field | Type | Description|
|:-----|:------|:------|
| `base58`  | string | The output identifier encoded with base58.    |
| `transactionID`   | string | The transaction identifier encoded with base58.     |
| `outputIndex`   | int | The index of an output.     |

#### Type `Consumers`

|Field | Type | Description|
|:-----|:------|:------|
| `transactionID`  | string | The transaction identifier encoded with base58.    |
| `valid`   | string | The boolean indicator if the transaction is valid.   |

<br />

## `/ledgerstate/outputs/:outputID/metadata`
Gets an output metadata for a given base58 encoded output ID.

### Parameters

| **Parameter**            | `outputID`      |
|--------------------------|----------------|
| **Required or Optional** | required       |
| **Description**          | The output ID encoded in base58. |
| **Type**                 | string         |

### Examples

#### cURL

```shell
curl http://localhost:8080/ledgerstate/outputs/:outputID/metadata \
-X GET \
-H 'Content-Type: application/json'

```

where `:outputID` is the ID of the output, e.g. 41GvDSQnd12e4nWnd2WzmdLmffruXqsE46jgeUbnB8s1QnK.

#### Client lib - `GetOutputMetadata()`
```Go
resp, err := goshimAPI.GetOutputMetadata("41GvDSQnd12e4nWnd2WzmdLmffruXqsE46jgeUbnB8s1QnK")
if err != nil {
    // return error
}
fmt.Printf("Metadata of an output %s:\n", resp.OutputID.Base58)
fmt.Println("branchID: ", resp.BranchID)
fmt.Println("first consumer: ", resp.FirstConsumer)
fmt.Println("number of consumers: ", resp.ConsumerCount)
fmt.Printf("finalized: %v, solid: %v\n", resp.Finalized, resp.Solid)
fmt.Println("solidification time: ",  time.Unix(resp.SolidificationTime, 0))
```
### Response examples
```json
{
    "outputID": {
        "base58": "41GvDSQnd12e4nWnd2WzmdLmffruXqsE46jgeUbnB8s1QnK",
        "transactionID": "9wr21zza46Y5QonKEHNQ6x8puA7Rbq5LAbsQZJCK1g1g",
        "outputIndex": 0
    },
    "branchID": "4uQeVj5tqViQh7yWWGStvkEG1Zmhx6uasJtWCJziofM",
    "solid": true,
    "solidificationTime": 1621889327,
    "consumerCount": 2,
    "firstConsumer": "HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV",
    "finalized": true
}
```

### Results

|Return field | Type | Description|
|:-----|:------|:------|
| `outputID`            | OutputID  | The output identifier encoded with base58.   |
| `branchID`            | string    | The identifier of the branch encoded with base58. |
| `solid`               | bool      | The boolean indicator if the message is solid. |
| `solidificationTime`  | int64     | The time of solidification of a message. |
| `consumerCount`       | int       | The number of consumers. |
| `firstConsumer`       | string    | The first consumer of the output. |
| `finalized`           | bool      | The boolean indicator if the transaction is finalized. |


#### Type `OutputID`

|Field | Type | Description|
|:-----|:------|:------|
| `base58`  | string | The output identifier encoded with base58.    |
| `transactionID`   | string | The transaction identifier encoded with base58.     |
| `outputIndex`   | int | The index of an output.     |

<br />

## `/ledgerstate/transactions/:transactionID`
Gets a transaction details for a given base58 encoded transaction ID.

### Parameters
| **Parameter**            | `transactionID`      |
|--------------------------|----------------|
| **Required or Optional** | required       |
| **Description**          | The transaction ID encoded in base58. |
| **Type**                 | string         |

### Examples

#### cURL

```shell
curl http://localhost:8080/ledgerstate/transactions/:transactionID \
-X GET \
-H 'Content-Type: application/json'
```

where `:transactionID` is the ID of the branch, e.g. HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV.

#### Client lib - `GetTransaction()`
```Go
resp, err := goshimAPI.GetTransaction("41GvDSQnd12e4nWnd2WzmdLmffruXqsE46jgeUbnB8s1QnK")
if err != nil {
    // return error
}
fmt.Println("transaction inputs:")
for _, input := range resp.Inputs {
    fmt.Println("inputID:", input.ReferencedOutputID.Base58)
}
fmt.Println("transaction outputs:")
for _, output := range resp.Outputs{
    fmt.Println("outputID:", output.OutputID.Base58)
    fmt.Println("output type:", output.Type)
}
fmt.Println("access mana pledgeID:", resp.AccessPledgeID)
fmt.Println("consensus mana pledgeID:", resp.ConsensusPledgeID)
```
### Response examples
```json
{
    "version": 0,
    "timestamp": 1621889348,
    "accessPledgeID": "DsHT39ZmwAGrKQe7F2rAjwHseUnJeY89gDPEH1FJxYdH",
    "consensusPledgeID": "DsHT39ZmwAGrKQe7F2rAjwHseUnJeY89gDPEH1FJxYdH",
    "inputs": [
        {
            "type": "UTXOInputType",
            "referencedOutputID": {
                "base58": "41GvDSQnd12e4nWnd2WzmdLmffruXqsE46jgeUbnB8s1QnK",
                "transactionID": "9wr21zza46Y5QonKEHNQ6x8puA7Rbq5LAbsQZJCK1g1g",
                "outputIndex": 0
            }
        }
    ],
    "outputs": [
        {
            "outputID": {
                "base58": "6gMWUCgJDozmyLeGzW3ibGFicEq2wbhsxgAw8rUVPvn9bj5",
                "transactionID": "HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV",
                "outputIndex": 0
            },
            "type": "SigLockedColoredOutputType",
            "output": {
                "balances": {
                    "11111111111111111111111111111111": 1000000
                },
                "address": "1HrUn1jWAjrMU58LLdFhfnWBwUKVdWjP5ojp7oCL9mVWs"
            }
        }
    ],
    "unlockBlocks": [
        {
            "type": "SignatureUnlockBlockType",
            "publicKey": "12vNcfgRHLSsobeqZFrjFRcVAmFQbDVniguPnEoxmkbG",
            "signature": "4isq3qzhY4MwbSeYM2NgRn5noWAyh5rqD12ruiTQ7P89TfXNecwHZ5nbpDc4UB7md1bkfM1xYtSh18FwLqK8HAC6"
        }
    ],
    "dataPayload": ""
}
```

### Results
|Return field | Type | Description|
|:-----|:------|:------|
| `version`         | uint8  | The version of the transaction essence.   |
| `timestamp`       | int64    | The issuing time of the transaction. |
| `accessPledgeID`  | string      | The node ID indicating to which node pledge the access mana. |
| `consensusPledgeID`| string     |  The node ID indicating to which node pledge the consensus mana. |
| `inputs`          | []Input       | The inputs of the transaction. |
| `outputs`         | []Output    | The outputs of the transaction. |
| `unlockBlocks`    | []UnlockBlock      | The unlock block containing signatures unlocking the inputs or references to previous unlock blocks. |
| `dataPayload` | []byte      | The raw data payload that can be attached to the transaction. |

#### Type `Input `
|Field | Type | Description|
|:-----|:------|:------|
| `Type`  | string | The type of input.   |
| `ReferencedOutputID`   | ReferencedOutputID | The output ID that is used as an input for the transaction.     |

#### Type `ReferencedOutputID`

|Field | Type | Description|
|:-----|:------|:------|
| `base58`  | string | The referenced output identifier encoded with base58.    |
| `transactionID`   | string | The transaction identifier encoded with base58.     |
| `outputIndex`   | int | The index of a referenced output.     |

#### Type `Output`

|Field | Type | Description|
|:-----|:------|:------|
| `outputID`  | OutputID | The identifier of an output.   |
| `outputType`   | string | The type of the output.  |
| `output` | string | An output raw message containing balances and corresponding addresses. |

#### Type `OutputID`

|Field | Type | Description|
|:-----|:------|:------|
| `base58`  | string | The output identifier encoded with base58.    |
| `transactionID`   | string | The transaction identifier encoded with base58.     |
| `outputIndex`   | int | The index of an output.     |

#### Type `UnlockBlock`
|Field | Type | Description|
|:-----|:------|:------|
| `type`  | string |  The unlock block type: signature or reference. |
| `referencedIndex`  | uint16 |  The reference index of an unlock block. |
| `signatureType`  | uint8 |  The unlock block signature type: ED25519 or BLS. |
| `publicKey`  | string |  The public key of a transaction owner. |
| `signature`  | string | The string representation of a signature encoded with base58 signed over a transaction essence. |

<br />

## `/ledgerstate/transactions/:transactionID/metadata`
Gets a transaction metadata for a given base58 encoded transaction ID.

### Parameters
| **Parameter**            | `transactionID`      |
|--------------------------|----------------|
| **Required or Optional** | required       |
| **Description**          | The transaction ID encoded in base58. |
| **Type**                 | string         |

### Examples

#### cURL

```shell
curl http://localhost:8080/ledgerstate/transactions/:transactionID/metadata \
-X GET \
-H 'Content-Type: application/json'
```

where `:transactionID` is the ID of the branch, e.g. HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV.

#### Client lib - `GetTransactionMetadata()`
```Go
resp, err := goshimAPI.GetTransactionMetadata("41GvDSQnd12e4nWnd2WzmdLmffruXqsE46jgeUbnB8s1QnK")
if err != nil {
    // return error
}
fmt.Println("transactionID:", resp.TransactionID)
fmt.Println("branchID:", resp.BranchID)
fmt.Printf("branch lazy booked: %v, solid: %v, finalized: %v\n", resp.LazyBooked, resp.Solid, resp.Finalized)
fmt.Println("solidification time:",  time.Unix(resp.SolidificationTime, 0))
```
### Response examples
```json
{
    "transactionID": "HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV",
    "branchID": "HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV",
    "solid": true,
    "solidificationTime": 1621889358,
    "finalized": true,
    "lazyBooked": false
}
```
### Results
|Return field | Type | Description|
|:-----|:------|:------|
| `transactionID`         | string  | The transaction identifier encoded with base58.    |
| `branchID`       | string    |  The branch identifier of the transaction. |
| `solid`  | bool      | The boolean indicator if the transaction is solid. |
| `solidificationTime`          | uint64      | The time of solidification of the transaction. |
| `finalized`         | bool    | The boolean indicator if the transaction is finalized. |
| `lazyBooked`    | bool      | The boolean indicator if the transaction is lazily booked.|

<br />

## `/ledgerstate/transactions/:transactionID/inclusionState`
Gets the fcob opinion associated with transaction based on a given base58 encoded transaction ID.

### Parameters
| **Parameter**            | `transactionID`      |
|--------------------------|----------------|
| **Required or Optional** | required       |
| **Description**          | The transaction ID encoded in base58. |
| **Type**                 | string         |

### Examples

#### cURL

```shell
curl http://localhost:8080/ledgerstate/transactions/:transactionID/inclusionState \
-X GET \
-H 'Content-Type: application/json'
```

where `:transactionID` is the ID of the branch, e.g. HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV.

#### Client lib - `GetTransactionInclusionState()`
```Go
resp, err := goshimAPI.GetTransactionInclusionState("41GvDSQnd12e4nWnd2WzmdLmffruXqsE46jgeUbnB8s1QnK")
if err != nil {
    // return error
}
fmt.Printf("Inclusion state of transaction %s:\n", resp.TransactionID)
fmt.Printf("conflicting: %v, confirmed: %v, rejected: %v, pending: %v\n", resp.Conflicting, resp.Confirmed, resp.Rejected, resp.Pending)
```
### Response examples
```json
{
    "transactionID": "HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV",
    "pending": false,
    "confirmed": true,
    "rejected": false,
    "conflicting": true
}
```

### Results
|Return field | Type | Description|
|:-----|:------|:------|
| `transactionID`         | string  | The transaction identifier encoded with base58.  |
| `pending`       | bool    | The boolean indicating if the transaction has not yet been confirmed nor rejected. |
| `confirmed`  | bool      | The boolean indicating if the transaction is confirmed. |
| `rejected`          | bool      | The boolean indicating if the transaction was rejected and is booked to the rejected branch. |
| `conflicting`         | bool    | The boolean indicating if the transaction is conflicting with some other transaction. |

<br />

## `/ledgerstate/transactions/:transactionID/consensus`
Gets the fcob opinion associated with transaction based on a given base58 encoded transaction ID.

### Parameters
| **Parameter**            | `transactionID`      |
|--------------------------|----------------|
| **Required or Optional** | required       |
| **Description**          | The transaction ID encoded in base58. |
| **Type**                 | string         |
### Examples

#### cURL

```shell
curl http://localhost:8080/ledgerstate/transactions/:transactionID/consensus \
-X GET \
-H 'Content-Type: application/json'
```

where `:transactionID` is the ID of the branch, e.g. HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV.

#### Client lib - `GetTransactionConsensusMetadata()`
```Go
resp, err := goshimAPI.GetTransactionConsensusMetadata("DNSN8GaCeep6CVuUV6KXAabXkL3bv4PUP4NkTNKoZMqS")
if err != nil {
    // return err
}
fmt.Printf("Consensus metadata for transaction %s:\n", resp.TransactionID)
fmt.Println("fcob opinion, is liked:", resp.Liked)
fmt.Println("level of knowledge:", resp.LoK)
fmt.Println("fcob timestamp:", time.Unix(resp.Timestamp, 0))
fmt.Println("fcob time 1:", time.Unix(resp.FCOBTime1, 0))
fmt.Println("fcob time 2:", time.Unix(resp.FCOBTime2, 0))
```

### Response examples
```json
{
    "transactionID": "HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV",
    "timestamp": 1621889358,
    "liked": true,
    "lok": "LevelOfKnowledge(Three)",
    "fcobTime1": 1621889360,
    "fcobTime2": 1621889362
}
```

### Results
|Return field | Type | Description|
|:-----|:------|:------|
| `transactionID`         | string  | The transaction identifier encoded with base58.  |
| `timestamp`       | uint64    | The timestamp of the arrival. |
| `lok`          | string      | The level of knowledge of the transaction - the degree of certainty of the associated opinion.|
| `liked`  | bool      | The fcob opinion for the transaction. |
| `fcobTime1`         | uint64    | The fcob opinion's execution time.|
| `fcobTime2`         | uint64    | The fcob opinion's locally finalized execution time.|

<br />

## `/ledgerstate/transactions/:transactionID/attachments`
Gets the list of messages IDs with attachments of the base58 encoded transaction ID.

### Parameters
| **Parameter**            | `transactionID`      |
|--------------------------|----------------|
| **Required or Optional** | required       |
| **Description**          | The transaction ID encoded in base58. |
| **Type**                 | string         |

### Examples

#### cURL

```shell
curl http://localhost:8080/ledgerstate/transactions/:transactionID/attachments \
-X GET \
-H 'Content-Type: application/json'
```

where `:transactionID` is the ID of the branch, e.g. HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV.

#### Client lib - `GetTransactionAttachments()`
```Go
resp, err := goshimAPI.GetTransactionAttachments("DNSN8GaCeep6CVuUV6KXAabXkL3bv4PUP4NkTNKoZMqS")
if err != nil {
    // return error
}
fmt.Printf("Messages IDs containing transaction %s:\n", resp.TransactionID)
for _, msgID := range resp.MessageIDs {
    fmt.Println(msgID)
}
```
### Response examples
```json
{
    "transactionID": "HuYUAwCeexmBePNXx5rNeJX1zUvUdUUs5LvmRmWe7HCV",
    "messageIDs": [
        "J1FQdMcticXiiuKMbjobq4zrYGHagk2mtTzkVwbqPgSq"
    ]
}
```

### Results
|Return field | Type | Description|
|:-----|:------|:------|
| `transactionID`   | string  | The transaction identifier encoded with base58.  |
| `messageIDs`       | []string    | The messages IDs that contains the requested transaction. |

<br />

## `/ledgerstate/transactions`
Sends transaction provided in form of a binary data, validates transaction before issuing the message payload. For more detail on how to prepare transaction bytes see the [tutorial]().

### Examples

#### Client lib - `PostTransaction()`
```GO
// prepare tx essence and signatures
...
// create transaction
tx := ledgerstate.NewTransaction(txEssence, ledgerstate.UnlockBlocks{unlockBlock})
resp, err := goshimAPI.PostTransaction(tx.Bytes())
if err != nil {
    // return error
}
fmt.Println("Transaction sent, txID: ", resp.TransactionID)
```

### Results
|Return field | Type | Description|
|:-----|:------|:------|
| `transactionID`   | string  | The transaction identifier encoded with base58.  |
| `Error`   | error  | The error returned if transaction was not processed correctly, otherwise is nil.  |

<br />

## `/ledgerstate/addresses/unspentOutputs`
Gets all unspent outputs for a list of addresses that were sent in the body message.  Returns the unspent outputs along with inclusion state and metadata for the wallet. 

### Request Body
```json
{
    "addresses": [
        "18LhfKUkWt4M9YR6Q3au4LT8wWCERwzHaqn153K78Eixp"
    ]
}
```

### Examples

#### cURL

```shell
curl http://localhost:8080/ledgerstate/addresses/unspentOutputs \
-X POST \
-H 'Content-Type: application/json'
--data-raw '{"addresses": ["18LhfKUkWt4M9YR6Q3au4LT8wWCERwzHaqn153K78Eixp"]}'
```

#### Client lib - `PostAddressUnspentOutputs()`
```Go
resp, err := goshimAPI.PostAddressUnspentOutputs([]string{"H36sZQkopfoEzP3WCMThSjUv5v9MLVYuaQ73tsKgVzXo"})
if err != nil {
    return
}
for _, outputs := range resp.UnspentOutputs {
    fmt.Println("address ID:", outputs.Address.Base58)
    fmt.Println("address type:", outputs.Address.Type)

    for _, output := range outputs.Outputs {
        fmt.Println("output ID:", output.Output.OutputID.Base58)
        fmt.Println("output type:", output.Output.Type)
    }
}
```

### Response examples

```json
{
    "unspentOutputs": [
        {
            "address": {
                "type": "AddressTypeED25519",
                "base58": "1Z4t5KEKU65fbeQCbNdztYTB1B4Cdxys1XRzTFrmvAf3"
            },
            "outputs": [
                {
                    "output": {
                        "outputID": {
                            "base58": "4eGoQWG7UDtBGK89vENQ5Ea1N1b8xF26VD2F8nigFqgyx5m",
                            "transactionID": "BqzgVk4yY9PDZuDro2mvT36U52ZYbJDfM41Xng3yWoQK",
                            "outputIndex": 0
                        },
                        "type": "SigLockedColoredOutputType",
                        "output": {
                            "balances": {
                                "11111111111111111111111111111111": 1000000
                            },
                            "address": "1Z4t5KEKU65fbeQCbNdztYTB1B4Cdxys1XRzTFrmvAf3"
                        }
                    },
                    "inclusionState": {
                        "confirmed": true,
                        "rejected": false,
                        "conflicting": false
                    },
                    "metadata": {
                        "timestamp": "2021-05-25T15:47:04.50470213+02:00"
                    }
                }
            ]
        }
    ]
}
```
### Results

|Return field | Type | Description|
|:-----|:------|:------|
| `unspentOutputs`   | WalletOutputsOnAddress  | Unspent outputs representation for wallet.  |

#### Type `WalletOutputsOnAddress`

|Return field | Type | Description|
|:-----|:------|:------|
| `Address`   | Address  | The address corresponding to the unspent output.  |
| `Outputs`   | []WalletOutput  | Unspent outputs representation for wallet.  |

#### Type `Address`

|Field | Type | Description|
|:-----|:------|:------|
| `type`  | string | The type of an address.   |
| `base58`| string   | The address encoded with base58.          |

#### Type `WalletOutput`

|Field | Type | Description|
|:-----|:------|:------|
| `output`  | Output | The unspent output.   |
| `inclusionState`| InclusionState   | The inclusion state of the transaction containing the output.  |
| `metadata`| WalletOutputMetadata   | The metadata of the output for the wallet lib.  |

#### Type `Output`

|Field | Type | Description|
|:-----|:------|:------|
| `outputID`  | OutputID | The identifier of an output.   |
| `outputType`   | string |  The type of the output.     |
| `output` | string | An outputs raw message containing balances and corresponding addresses |

#### Type `OutputID`

|Field | Type | Description|
|:-----|:------|:------|
| `base58`  | string | The output identifier encoded with base58.    |
| `transactionID`   | string | The transaction identifier encoded with base58.     |
| `outputIndex`   | int | The index of an output.     |

#### Type `InclusionState`

|Field | Type | Description|
|:-----|:------|:------|
| `confirmed`  | bool |  The boolean indicating if the transaction containing the output is confirmed. |
| `rejected`   | bool | The boolean indicating if the transaction that contains the output was rejected and is booked to the rejected branch.   |
| `conflicting`   | bool | The boolean indicating if the output is in conflicting transaction.|

#### Type `WalletOutputMetadata`

|Field | Type | Description|
|:-----|:------|:------|
| `timestamp`  | time.Time | The timestamp of the transaction containing the output.    |
