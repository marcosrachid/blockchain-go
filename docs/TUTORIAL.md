# Tutorial: Primeiros Passos com o Blockchain

Este tutorial mostrará como usar o blockchain passo a passo.

## Passo 1: Compilar o Projeto

```bash
make build
# ou
go build -o blockchain
```

## Passo 2: Criar Carteiras

Primeiro, vamos criar duas carteiras (uma para Alice e outra para Bob):

```bash
./blockchain createwallet
# Saída: New address is: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa

./blockchain createwallet  
# Saída: New address is: 1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2
```

**Importante**: Anote esses endereços! Vamos chamá-los de:
- `ALICE_ADDRESS`: primeiro endereço criado
- `BOB_ADDRESS`: segundo endereço criado

## Passo 3: Listar Endereços

Você pode ver todos os endereços criados:

```bash
./blockchain listaddresses
```

## Passo 4: Criar o Blockchain

Crie o blockchain enviando a recompensa do bloco gênesis para Alice:

```bash
./blockchain createblockchain -address ALICE_ADDRESS
```

Você verá:
```
No existing blockchain found
Genesis created
Done! There are 1 transactions in the UTXO set.
Finished!
```

## Passo 5: Verificar Saldo

Alice deve ter 50 moedas (recompensa do bloco gênesis):

```bash
./blockchain getbalance -address ALICE_ADDRESS
# Saída: Balance of ALICE_ADDRESS: 50
```

Bob ainda não tem saldo:

```bash
./blockchain getbalance -address BOB_ADDRESS
# Saída: Balance of BOB_ADDRESS: 0
```

## Passo 6: Enviar Transação

Alice envia 10 moedas para Bob:

```bash
./blockchain send -from ALICE_ADDRESS -to BOB_ADDRESS -amount 10
```

Durante a mineração, você verá o Proof of Work em ação:
```
Hash: 0000abc123..., Nonce: 12345
Success!
```

## Passo 7: Verificar Novos Saldos

Verifique o saldo de Alice (deve ter 40 + 50 da recompensa de mineração):

```bash
./blockchain getbalance -address ALICE_ADDRESS
# Saída: Balance of ALICE_ADDRESS: 90
```

Explicação:
- Alice tinha 50
- Enviou 10 para Bob
- Recebeu 50 de recompensa por minerar o bloco
- Total: 50 - 10 + 50 = 90

Verifique o saldo de Bob:

```bash
./blockchain getbalance -address BOB_ADDRESS
# Saída: Balance of BOB_ADDRESS: 10
```

## Passo 8: Visualizar a Blockchain

```bash
./blockchain printchain
```

Você verá algo como:

```
============ Block 0000abc123def456... ============
Height: 1
Prev. hash: 0000xyz789...
PoW: true
--- Transaction abc123...:
     Input 0:
       TXID:     def456...
       Out:       0
       Signature: 789abc...
       PubKey:    123def...
     Output 0:
       Value:  10
       Script: 456789...
     Output 1:
       Value:  40
       Script: abc123...


============ Block 0000def456abc789... ============
Height: 0
Prev. hash: 
PoW: true
--- Transaction (Genesis):
     Input 0:
       TXID:     
       Out:       -1
     Output 0:
       Value:  50
       Script: xyz123...
```

## Entendendo os Conceitos

### 1. **Recompensa de Mineração**
Cada vez que você envia uma transação, um bloco é minerado e você recebe 50 moedas como recompensa.

### 2. **Troco**
Se Alice envia 10 moedas mas tem um UTXO de 50, o sistema cria:
- Output de 10 para Bob
- Output de 40 de troco para Alice

### 3. **Proof of Work**
O minerador precisa encontrar um nonce que faça o hash do bloco começar com zeros. Quanto maior a dificuldade, mais zeros são necessários.

### 4. **UTXOs**
Cada transação consome outputs de transações anteriores (inputs) e cria novos outputs. Um output só pode ser gasto uma vez.

## Cenário Completo de Teste

Aqui está um script completo para testar:

```bash
# 1. Compilar
make build

# 2. Criar carteiras
echo "Criando carteira para Alice..."
ALICE=$(./blockchain createwallet | grep "New address is:" | cut -d' ' -f4)
echo "Alice: $ALICE"

echo "Criando carteira para Bob..."
BOB=$(./blockchain createwallet | grep "New address is:" | cut -d' ' -f4)
echo "Bob: $BOB"

echo "Criando carteira para Charlie..."
CHARLIE=$(./blockchain createwallet | grep "New address is:" | cut -d' ' -f4)
echo "Charlie: $CHARLIE"

# 3. Criar blockchain
echo "Criando blockchain..."
./blockchain createblockchain -address $ALICE

# 4. Verificar saldos iniciais
echo "Saldos iniciais:"
./blockchain getbalance -address $ALICE
./blockchain getbalance -address $BOB
./blockchain getbalance -address $CHARLIE

# 5. Alice envia 10 para Bob
echo "Alice envia 10 para Bob..."
./blockchain send -from $ALICE -to $BOB -amount 10

# 6. Alice envia 20 para Charlie
echo "Alice envia 20 para Charlie..."
./blockchain send -from $ALICE -to $CHARLIE -amount 20

# 7. Bob envia 5 para Charlie
echo "Bob envia 5 para Charlie..."
./blockchain send -from $BOB -to $CHARLIE -amount 5

# 8. Verificar saldos finais
echo "Saldos finais:"
./blockchain getbalance -address $ALICE
./blockchain getbalance -address $BOB
./blockchain getbalance -address $CHARLIE

# 9. Imprimir blockchain
echo "Blockchain completa:"
./blockchain printchain
```

## Comandos Úteis

### Limpar tudo e começar de novo:
```bash
make clean
make build
```

### Verificar quantas transações no UTXO set:
```bash
./blockchain reindexutxo
```

### Usar o Makefile:
```bash
# Criar carteira
make wallet

# Criar blockchain
make blockchain ADDRESS=$ALICE

# Enviar transação
make send FROM=$ALICE TO=$BOB AMOUNT=10

# Ver saldo
make balance ADDRESS=$ALICE
```

## Problemas Comuns

### 1. "No existing blockchain found"
Você precisa criar o blockchain primeiro:
```bash
./blockchain createblockchain -address SEU_ENDERECO
```

### 2. "Not enough funds"
Você não tem moedas suficientes. Verifique seu saldo primeiro.

### 3. "Address is not Valid"
Verifique se você está usando um endereço válido criado pelo `createwallet`.

### 4. Banco de dados travado
Se o programa foi interrompido abruptamente:
```bash
rm -rf ./tmp
# Depois crie o blockchain novamente
```

## Explorando o Código

### Ver como funciona o Proof of Work:
```go
// blockchain/proof.go
func (pow *ProofOfWork) Run() (int, []byte) {
    // Loop até encontrar hash válido
    for nonce < math.MaxInt64 {
        data := pow.InitData(nonce)
        hash = sha256.Sum256(data)
        // Verifica se hash é menor que target
        if intHash.Cmp(pow.Target) == -1 {
            break
        }
        nonce++
    }
}
```

### Ver como transações são criadas:
```go
// blockchain/transaction.go
func NewTransaction(from, to string, amount int, chain *Blockchain) *Transaction {
    // Encontra UTXOs gastáveis
    acc, validOutputs := chain.FindSpendableOutputs(pubKeyHash, amount)
    // Cria inputs dos UTXOs
    // Cria outputs (um para destino, outro para troco)
    // Assina a transação
}
```

### Ver como Merkle Tree funciona:
```go
// blockchain/merkle.go
func NewMerkleTree(data [][]byte) *MerkleTree {
    // Cria folhas da árvore
    // Constrói árvore de baixo para cima
    // Retorna raiz
}
```

## Próximos Experimentos

1. **Modificar a dificuldade**: Edite `Difficulty` em `blockchain/proof.go` (valores menores = mais fácil)

2. **Modificar recompensa**: Edite `subsidy` em `blockchain/transaction.go`

3. **Adicionar mais funcionalidades**:
   - Sistema de taxas de transação
   - Limitar tamanho do bloco
   - Adicionar timestamp nas transações
   - Criar sistema de "notas" nas transações

4. **Estudar o código**: Leia os arquivos na ordem:
   - `block.go` → estrutura básica
   - `proof.go` → mineração
   - `transaction.go` → transações
   - `wallet.go` → carteiras
   - `blockchain.go` → blockchain completo
   - `utxo.go` → otimizações
   - `merkle.go` → estrutura de dados

---

Divirta-se explorando o blockchain! 🚀

