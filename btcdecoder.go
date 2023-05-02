package main

import (
    "encoding/hex"
    "encoding/binary"
    "fmt"
)

type TransactionInput struct {
    PreviousTxHash [32]byte
    OutputIndex    uint32
    ScriptLength   uint8
    SignatureScript []byte
    Sequence       uint32
}

type TransactionOutput struct {
    Value           int64
    ScriptLength    uint8
    PublicKeyScript []byte
}

type Transaction struct {
    Version      uint32
    InputCount   uint8
    Inputs       []TransactionInput
    OutputCount  uint8
    Outputs      []TransactionOutput
    LockTime     uint32
}

func DecodeTransaction(hexTx string) (*Transaction, error) {
    var tx Transaction

    txBytes, err := hex.DecodeString(hexTx)
    if err != nil {
        return nil, err
    }

    offset := 0

    // Read the transaction version
    tx.Version = binary.LittleEndian.Uint32(txBytes[offset:])
    offset += 4

    // Read the number of inputs
    tx.InputCount = uint8(txBytes[offset])
    offset++

    // Decode each input
    for i := 0; i < int(tx.InputCount); i++ {
        var input TransactionInput

        // Read the previous transaction hash
        copy(input.PreviousTxHash[:], txBytes[offset:offset+32])
        offset += 32

        // Read the output index
        input.OutputIndex = binary.LittleEndian.Uint32(txBytes[offset:])
        offset += 4

        // Read the length of the signature script
        input.ScriptLength = uint8(txBytes[offset])
        offset++

        // Read the signature script
        input.SignatureScript = make([]byte, input.ScriptLength)
        copy(input.SignatureScript, txBytes[offset:offset+int(input.ScriptLength)])
        offset += int(input.ScriptLength)

        // Read the sequence number
        input.Sequence = binary.LittleEndian.Uint32(txBytes[offset:])
        offset += 4

        tx.Inputs = append(tx.Inputs, input)
    }

    // Read the number of outputs
    tx.OutputCount = uint8(txBytes[offset])
    offset++

    // Decode each output
    for i := 0; i < int(tx.OutputCount); i++ {
        var output TransactionOutput

        // Read the output value
        output.Value = int64(binary.LittleEndian.Uint64(txBytes[offset:]))
        offset += 8

        // Read the length of the public key script
        output.ScriptLength = uint8(txBytes[offset])
        offset++

        // Read the public key script
        output.PublicKeyScript = make([]byte, output.ScriptLength)
        copy(output.PublicKeyScript, txBytes[offset:offset+int(output.ScriptLength)])
        offset += int(output.ScriptLength)

        tx.Outputs = append(tx.Outputs, output)
    }

    // Read the lock time
    tx.LockTime = binary.LittleEndian.Uint32(txBytes[offset:])

    return &tx, nil
}
