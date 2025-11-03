# Blockchain em Go - Similar ao Bitcoin

Um projeto educacional de blockchain implementado em Go, inspirado no protocolo do Bitcoin.

## 📚 Características Implementadas

Este projeto implementa os principais conceitos do Bitcoin:

### 1. **Proof of Work (PoW)**
- Algoritmo de consenso similar ao Bitcoin
- Dificuldade ajustável
- Mineração de blocos com nonce

### 2. **Sistema de Transações**
- Transações com múltiplos inputs e outputs
- Transações Coinbase (recompensa de mineração)
- Verificação e assinatura digital de transações

### 3. **UTXOs (Unspent Transaction Outputs)**
- Modelo UTXO similar ao Bitcoin
- Cache de UTXOs para performance
- Sistema de tracking de outputs não gastos

### 4. **Criptografia ECDSA**
- Geração de pares de chaves (pública/privada)
- Assinatura digital de transações usando ECDSA
- Curva elíptica P256

### 5. **Carteiras (Wallets)**
- Geração de endereços Bitcoin-like
- Codificação Base58 (alfabeto Bitcoin)
- Hash de chave pública (SHA256 + RIPEMD160)
- Checksum para validação de endereços

### 6. **Merkle Tree**
- Estrutura de dados para verificar transações
- Hash eficiente de todas as transações do bloco
- Usado no cabeçalho do bloco

### 7. **Persistência**
- Banco de dados LevelDB (suporta acesso concorrente de leitura/escrita)
- Serialização/deserialização de blocos
- Iterador para percorrer a blockchain

### 8. **CLI (Interface de Linha de Comando)**
- Comandos para interagir com o blockchain
- Criação de carteiras
- Envio de transações
- Consulta de saldos

## 🏗️ Estrutura do Projeto

```
blockchain-go/
├── blockchain/
│   ├── base58.go          # Codificação Base58 (Bitcoin)
│   ├── block.go           # Estrutura de Block com transações
│   ├── blockchain.go      # Blockchain principal com UTXO
│   ├── merkle.go          # Implementação de Merkle Tree
│   ├── proof.go           # Proof of Work
│   ├── transaction.go     # Sistema de transações
│   ├── utxo.go            # Sistema UTXO
│   ├── utils.go           # Funções utilitárias
│   └── wallet.go          # Sistema de carteiras
├── cli/
│   └── cli.go             # Interface de linha de comando
├── main.go                # Ponto de entrada
├── go.mod                 # Dependências
└── README.md
```

## 🚀 Como Usar

### Compilar

```bash
go build -o blockchain-app
# ou
make build
```

### Criar uma Carteira

```bash
./blockchain-app createwallet
# ou
make wallet
```

Isso gerará um novo endereço Bitcoin-like, exemplo:
```
New address is: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
```

### Listar Endereços

```bash
./blockchain-app listaddresses
# ou
make list
```

### Criar o Blockchain

Crie o blockchain e envie a recompensa do bloco gênesis para um endereço:

```bash
./blockchain-app createblockchain -address SEU_ENDERECO
# ou
make blockchain ADDRESS=SEU_ENDERECO
```

### Verificar Saldo

```bash
./blockchain-app getbalance -address SEU_ENDERECO
# ou
make balance ADDRESS=SEU_ENDERECO
```

### Enviar Transação

```bash
./blockchain-app send -from ENDERECO_ORIGEM -to ENDERECO_DESTINO -amount 10
# ou
make send FROM=ENDERECO_ORIGEM TO=ENDERECO_DESTINO AMOUNT=10
```

### Visualizar a Blockchain

```bash
./blockchain-app printchain
# ou
make print
```

### Reindexar UTXOs

```bash
./blockchain-app reindexutxo
# ou
make reindex
```

## 📖 Conceitos do Bitcoin Implementados

### 1. Proof of Work
O algoritmo de consenso que garante segurança através de trabalho computacional:
- Mineradores devem encontrar um hash que atenda à dificuldade estabelecida
- O hash deve ter um certo número de zeros à esquerda
- Similar ao SHA256(SHA256()) do Bitcoin

### 2. Transações
Estrutura similar ao Bitcoin:
- **Inputs**: Referências a outputs de transações anteriores
- **Outputs**: Novos destinos para as moedas com valores específicos
- **Coinbase**: Transação especial de recompensa para o minerador

### 3. UTXO (Unspent Transaction Output)
- Modelo de contabilidade do Bitcoin
- Cada output só pode ser gasto uma vez
- Sistema de tracking de outputs não gastos para eficiência

### 4. Criptografia
- **ECDSA**: Assinatura digital de transações
- **SHA256**: Hashing de blocos e transações
- **RIPEMD160**: Hashing de chaves públicas
- **Base58**: Codificação de endereços (evita caracteres ambíguos)

### 5. Merkle Tree
- Estrutura de dados que permite verificação eficiente de transações
- Raiz da árvore incluída no cabeçalho do bloco
- Permite SPV (Simplified Payment Verification)

### 6. Estrutura do Bloco
```go
type Block struct {
    Timestamp    int64           // Quando o bloco foi minerado
    Hash         []byte          // Hash do bloco
    Transactions []*Transaction  // Transações no bloco
    PrevHash     []byte          // Hash do bloco anterior
    Nonce        int             // Nonce para PoW
    Height       int             // Altura na blockchain
}
```

### 7. Carteiras e Endereços
Processo de geração de endereço similar ao Bitcoin:
1. Gera par de chaves ECDSA
2. SHA256 da chave pública
3. RIPEMD160 do resultado
4. Adiciona byte de versão
5. Calcula checksum (SHA256(SHA256()))
6. Codifica em Base58

## 🔍 Diferenças em Relação ao Bitcoin Real

Este é um projeto educacional. Algumas diferenças em relação ao Bitcoin real:

1. **Dificuldade Fixa**: Bitcoin ajusta dificuldade a cada 2016 blocos
2. **Recompensa Fixa**: Bitcoin reduz recompensa pela metade a cada 210.000 blocos (halving)
3. **Rede P2P**: Não implementado (Bitcoin tem protocolo de rede completo)
4. **Scripts**: Bitcoin usa linguagem Script para condições de gasto
5. **Mempool**: Pool de transações pendentes não implementado
6. **SPV**: Simplified Payment Verification não implementado
7. **Segregated Witness**: Não implementado
8. **Lightning Network**: Não implementado

## 🛠️ Tecnologias Utilizadas

- **Go 1.24+**: Linguagem de programação
- **LevelDB**: Banco de dados key-value para persistência (com suporte a acesso concorrente)
- **crypto/ecdsa**: Criptografia de curva elíptica
- **crypto/sha256**: Função hash SHA-256
- **golang.org/x/crypto/ripemd160**: Hash RIPEMD-160

## 📚 Recursos para Aprendizado

Para entender melhor o Bitcoin:

1. [Bitcoin Whitepaper - Satoshi Nakamoto](https://bitcoin.org/bitcoin.pdf)
2. [Mastering Bitcoin - Andreas Antonopoulos](https://github.com/bitcoinbook/bitcoinbook)
3. [Bitcoin Developer Guide](https://bitcoin.org/en/developer-guide)
4. [Learn Me a Bitcoin](https://learnmeabitcoin.com/)

## 🤝 Contribuindo

Este é um projeto educacional. Sinta-se livre para:
- Fazer fork do projeto
- Adicionar novas funcionalidades
- Melhorar a documentação
- Reportar issues

## ⚠️ Aviso

Este projeto foi criado apenas para fins educacionais e não deve ser usado em produção. Não é adequado para armazenar valores reais.

## 📄 Licença

Este projeto é de código aberto e está disponível para uso educacional.

## 🎯 Próximos Passos Sugeridos

Para expandir o projeto e torná-lo ainda mais similar ao Bitcoin:

1. **Ajuste de Dificuldade Dinâmico**: Implementar ajuste automático baseado no tempo de mineração
2. **Halving de Recompensa**: Reduzir recompensa pela metade em intervalos específicos
3. **Rede P2P**: Adicionar capacidade de comunicação entre nós
4. **Mempool**: Pool de transações pendentes aguardando mineração
5. **Script System**: Sistema de scripts para condições de gasto mais complexas
6. **SegWit**: Implementar Segregated Witness
7. **Interface Web**: Criar uma interface web para visualizar o blockchain
8. **Testes Unitários**: Adicionar cobertura de testes completa
9. **Métricas**: Adicionar estatísticas (hashrate, tamanho do blockchain, etc.)
10. **API REST**: Criar API para integração com outras aplicações

---

**Desenvolvido com 💙 para aprendizado de Blockchain e Bitcoin**

