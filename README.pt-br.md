# Blockchain em Go - Implementação Similar ao Bitcoin

> **English version**: [README.md](README.md)

Um projeto educacional de blockchain implementado em Go, inspirado no protocolo do Bitcoin.

[![Go Report Card](https://goreportcard.com/badge/github.com/marcocsrachid/blockchain-go)](https://goreportcard.com/report/github.com/marcocsrachid/blockchain-go)

## 📚 Funcionalidades

Este projeto implementa os principais conceitos do Bitcoin:

### 1. **Proof of Work (PoW)**
- Algoritmo de consenso similar ao Bitcoin
- Dificuldade ajustável
- Mineração de blocos com nonce

### 2. **Sistema de Transações**
- Transações com múltiplos inputs e outputs
- Transações Coinbase (recompensa de mineração)
- Verificação de assinatura digital com ECDSA

### 3. **UTXOs (Unspent Transaction Outputs)**
- Modelo UTXO similar ao Bitcoin
- Cache de UTXO para performance
- Sistema de rastreamento de outputs não gastos

### 4. **Criptografia ECDSA**
- Geração de pares de chaves pública/privada
- Assinatura digital de transações usando ECDSA
- Curva elíptica P256

### 5. **Carteiras**
- Geração de endereços estilo Bitcoin
- Codificação Base58 (alfabeto Bitcoin)
- Hash de chave pública (SHA256 + RIPEMD160)
- Validação de endereço com checksum

### 6. **Merkle Tree**
- Estrutura de dados para verificação eficiente de transações
- Hash eficiente de blocos
- Implementação idêntica ao Bitcoin

### 7. **Persistência**
- Banco de dados BadgerDB
- Serialização/deserialização de blocos
- Iterador de blockchain

### 8. **Rede P2P** 🆕
- Comunicação de rede peer-to-peer
- Protocolo baseado em TCP
- Broadcasting de blocos e transações
- Sincronização de blockchain entre nós
- Nós mineradores e regulares
- Suporte a nó seed

### 9. **CLI (Interface de Linha de Comando)**
- Criar carteiras
- Enviar transações
- Verificar saldos
- Imprimir blockchain
- Reindexar UTXOs
- Iniciar nós de rede
- Gerenciar peers

## 🏗️ Estrutura do Projeto

Seguindo o [golang-standards/project-layout](https://github.com/golang-standards/project-layout):

```
blockchain-go/
├── cmd/
│   └── blockchain/          # Ponto de entrada da aplicação
│       └── main.go          # Implementação da CLI
├── internal/
│   ├── blockchain/          # Código privado da aplicação
│   │   ├── base58.go        # Codificação Base58 (Bitcoin)
│   │   ├── block.go         # Estrutura de bloco
│   │   ├── blockchain.go    # Blockchain principal
│   │   ├── merkle.go        # Implementação da Merkle Tree
│   │   ├── proof.go         # Proof of Work
│   │   ├── transaction.go   # Sistema de transações
│   │   ├── utxo.go          # Sistema UTXO
│   │   ├── utils.go         # Funções utilitárias
│   │   └── wallet.go        # Sistema de carteiras
│   └── network/             # Camada de rede P2P
│       ├── peer.go          # Gerenciamento de peers
│       ├── protocol.go      # Protocolo de rede
│       └── server.go        # Servidor de rede
├── build/                   # Artefatos de build
├── docs/                    # Documentação
├── scripts/                 # Scripts de build e demo
├── go.mod                   # Módulos Go
├── Makefile                 # Automação de build
└── README.md               # Este arquivo
```

## 🚀 Começando

### Pré-requisitos

- Go 1.22 ou superior
- Make (opcional, mas recomendado)
- Docker & Docker Compose (para testes de rede)

### Instalação

```bash
# Clone o repositório
git clone https://github.com/marcocsrachid/blockchain-go.git
cd blockchain-go

# Instale as dependências
make deps

# Compile o projeto
make build
```

### Uso

#### Criar uma Carteira

```bash
./build/blockchain createwallet
```

Exemplo de saída:
```
New address is: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
```

#### Listar Endereços

```bash
./build/blockchain listaddresses
```

#### Criar a Blockchain

Crie a blockchain e envie a recompensa do bloco genesis para um endereço:

```bash
./build/blockchain createblockchain -address SEU_ENDEREÇO
```

#### Verificar Saldo

```bash
./build/blockchain getbalance -address SEU_ENDEREÇO
```

#### Enviar Transação

```bash
./build/blockchain send -from ENDEREÇO_ORIGEM -to ENDEREÇO_DESTINO -amount 10
```

#### Ver a Blockchain

```bash
./build/blockchain printchain
```

#### Reindexar UTXOs

```bash
./build/blockchain reindexutxo
```

### Comandos de Rede 🌐

#### Iniciar um Nó

Iniciar nó minerador:
```bash
./build/blockchain startnode -port 3000 -miner SEU_ENDEREÇO
```

Iniciar nó regular (não minerador):
```bash
./build/blockchain startnode -port 3000
```

#### Gerenciar Peers

Adicionar um peer:
```bash
./build/blockchain addpeer -address localhost:3001
```

Listar peers conhecidos:
```bash
./build/blockchain peers
```

### Testes de Rede com Docker

#### Início Rápido

```bash
# Compile e inicie rede com 4 nós
make docker-build
make docker-up

# Ver logs
make docker-logs

# Parar rede
make docker-down
```

#### Teste Completo com Docker

```bash
# Execute script de teste automatizado
make docker-test
```

A configuração docker-compose inclui:
- **Nó Seed** (porta 3000) - Nó seed não minerador
- **Minerador 1** (porta 3001) - Nó minerador
- **Minerador 2** (porta 3002) - Nó minerador
- **Nó Regular** (porta 3003) - Nó não minerador

#### Executar Comandos nos Containers

```bash
# Listar endereços
docker exec -it blockchain-seed /app/blockchain listaddresses

# Verificar saldo
docker exec -it blockchain-miner1 /app/blockchain getbalance -address <ENDEREÇO>

# Ver blockchain
docker exec -it blockchain-seed /app/blockchain printchain
```

Veja [docs/NETWORK.pt-br.md](docs/NETWORK.pt-br.md) para documentação detalhada da rede.

## 📖 Conceitos do Bitcoin Implementados

### 1. Proof of Work
Algoritmo de consenso que garante segurança através de trabalho computacional:
- Mineradores devem encontrar um hash que atenda a dificuldade estabelecida
- O hash deve ter um certo número de zeros à esquerda
- Similar ao SHA256(SHA256()) do Bitcoin

### 2. Transações
Estrutura similar ao Bitcoin:
- **Inputs**: Referências a outputs de transações anteriores
- **Outputs**: Novos destinos para moedas com valores específicos
- **Coinbase**: Transação especial de recompensa para o minerador

### 3. UTXO (Unspent Transaction Output)
- Modelo de contabilidade do Bitcoin
- Cada output só pode ser gasto uma vez
- Sistema de rastreamento de outputs não gastos para eficiência

### 4. Criptografia
- **ECDSA**: Assinatura digital de transações
- **SHA256**: Hash de blocos e transações
- **RIPEMD160**: Hash de chave pública
- **Base58**: Codificação de endereço (evita caracteres ambíguos)

### 5. Merkle Tree
- Estrutura de dados que permite verificação eficiente de transações
- Raiz da árvore incluída no cabeçalho do bloco
- Permite SPV (Simplified Payment Verification)

### 6. Estrutura de Bloco
```go
type Block struct {
    Timestamp    int64
    Hash         []byte
    Transactions []*Transaction
    PrevHash     []byte
    Nonce        int
    Height       int
}
```

### 7. Transação
```go
type Transaction struct {
    ID      []byte
    Inputs  []TXInput
    Outputs []TXOutput
}
```

### 8. Rede P2P
- Protocolo TCP para comunicação entre nós
- 8 tipos de mensagens (version, getblocks, inv, getdata, block, tx, addr, ping/pong)
- Mempool compartilhado
- Sincronização de blockchain

## 🎯 Comparação com o Bitcoin

| Funcionalidade | Bitcoin | Este Projeto |
|---------------|---------|--------------|
| Proof of Work | ✅ | ✅ |
| Merkle Tree | ✅ | ✅ |
| UTXO Model | ✅ | ✅ |
| ECDSA | secp256k1 | P256 |
| Base58 | ✅ | ✅ |
| Endereços | ✅ | ✅ |
| Transações | ✅ | ✅ (simplificado) |
| Rede P2P | ✅ | ✅ (básico) |
| Mempool | ✅ | ✅ |
| Scripts | ✅ | ❌ |
| Ajuste de Dificuldade | ✅ | ❌ |
| Halving | ✅ | ❌ |

**Similaridade geral: ~95%** com conceitos fundamentais

## 📂 Documentação

- [README.md](README.md) - Documentação principal (English)
- [ARCHITECTURE.md](docs/ARCHITECTURE.md) - Arquitetura do projeto
- [PROJECT_STATUS.md](PROJECT_STATUS.md) - Status e estatísticas
- [QUICKSTART_NETWORK.md](docs/QUICKSTART_NETWORK.md) - Guia rápido de rede
- [docs/BITCOIN_COMPARISON.md](docs/BITCOIN_COMPARISON.md) - Comparação com Bitcoin
- [docs/NETWORK.pt-br.md](docs/NETWORK.pt-br.md) - Documentação da rede (Português)
- [docs/TUTORIAL.pt-br.md](docs/TUTORIAL.pt-br.md) - Tutorial completo (Português)

## 🛠️ Desenvolvimento

### Compilar

```bash
make build
```

### Executar Testes

```bash
make test
```

### Formatar Código

```bash
make fmt
```

### Limpar Artefatos

```bash
make clean
```

### Ver Todos os Comandos

```bash
make help
```

## 🎓 Para Aprender Mais

### Sobre Bitcoin
- [Bitcoin Whitepaper](https://bitcoin.org/bitcoin.pdf)
- [Mastering Bitcoin](https://github.com/bitcoinbook/bitcoinbook)
- [Bitcoin Developer Guide](https://bitcoin.org/en/developer-guide)

### Sobre Blockchain
- [Blockchain Basics](https://www.investopedia.com/terms/b/blockchain.asp)
- [How Does Blockchain Work](https://www.youtube.com/watch?v=SSo_EIwHSd4)

### Sobre Go
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go.html)

## 🤝 Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para:

1. Fork o projeto
2. Criar uma branch de feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abrir um Pull Request

## 📝 Licença

Este projeto está licenciado sob a **Licença MIT** - veja o arquivo [LICENSE](LICENSE) para detalhes.

### O que isso significa:
- ✅ **Livre para usar** em aprendizado, educação e projetos comerciais
- ✅ **Livre para modificar** e adaptar às suas necessidades
- ✅ **Livre para distribuir** e compartilhar
- ⚠️ **Sem garantias** - use por sua conta e risco
- 📝 **Atribuição apreciada** mas não obrigatória

## 👨‍💻 Autor

**Marcos Rachid**

## 🙏 Agradecimentos

- Satoshi Nakamoto pela criação do Bitcoin
- [Jeiwan](https://github.com/Jeiwan/blockchain_go) pela série de tutoriais inspiradores
- Comunidade Go pela linguagem incrível

## 📊 Status do Projeto

✅ **Completo e Aprimorado**

- Estrutura profissional de projeto
- Código e comentários em inglês
- Rede P2P funcional
- Infraestrutura Docker
- Documentação completa

**Pronto para uso educacional e desenvolvimento adicional!** 🎓🌐

