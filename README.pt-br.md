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
- Banco de dados LevelDB (suporta acesso concorrente de leitura/escrita)
- Serialização/deserialização de blocos
- Iterador de blockchain

### 8. **Rede P2P** 🆕
- Comunicação de rede peer-to-peer
- Protocolo baseado em TCP
- Broadcasting de blocos e transações
- Sincronização de blockchain entre nós
- Nós mineradores e regulares
- Suporte a nó seed

### 9. **API REST HTTP**

- Criar carteiras (`POST /api/createwallet`)
- Enviar transações (`POST /api/send`)
- Verificar saldos (`GET /api/balance/:address`)
- Info da rede (`GET /api/networkinfo`)
- Listar endereços (`GET /api/addresses`)
- Ver último bloco (`GET /api/lastblock`)
- Health check (`GET /health`)

### 10. **CLI (Interface de Linha de Comando)**

- Iniciar nós de rede (`startnode`)
- Criar blockchain (`createblockchain`)
- Gerenciamento básico de carteiras (`createwallet`, `listaddresses`)

## 🏗️ Estrutura do Projeto

Seguindo o [golang-standards/project-layout](https://github.com/golang-standards/project-layout):

```
blockchain-go/
├── cmd/
│   └── blockchain/          # Ponto de entrada da aplicação
│       └── main.go          # Inicialização e comandos básicos
├── internal/
│   ├── api/                 # Servidor HTTP API
│   │   └── server.go        # Endpoints REST (balance, send, network info, etc.)
│   ├── blockchain/          # Lógica core da blockchain
│   │   ├── base58.go        # Codificação Base58 (estilo Bitcoin)
│   │   ├── block.go         # Estrutura de bloco com PoW e transações
│   │   ├── blockchain.go    # Blockchain com persistência (LevelDB)
│   │   ├── config.go        # Constantes de configuração (dificuldade, rewards, etc.)
│   │   ├── merkle.go        # Merkle Tree para hash de transações
│   │   ├── proof.go         # Algoritmo Proof of Work
│   │   ├── transaction.go   # Sistema de transações com assinaturas ECDSA
│   │   ├── utxo.go          # Gerenciamento do conjunto UTXO
│   │   ├── utils.go         # Funções utilitárias
│   │   └── wallet.go        # Gerenciamento de carteiras e endereços
│   └── network/             # Camada de rede P2P
│       ├── peer.go          # Gerenciamento de conexões de peers
│       ├── protocol.go      # Mensagens do protocolo de rede
│       └── server.go        # Servidor P2P, mempool, coordenação de mineração
├── build/                   # Binários compilados
├── docs/                    # Documentação detalhada
│   ├── ARCHITECTURE.md      # Arquitetura do sistema
│   ├── BITCOIN_COMPARISON.md # Comparação com Bitcoin
│   ├── HALVING_AND_SUPPLY.md # Modelo econômico e supply
│   ├── MINING.md            # Mecânicas de mineração
│   ├── NETWORK.md           # Detalhes do protocolo de rede
│   └── ...                  # Versões em português (*.pt-br.md)
├── scripts/                 # Scripts utilitários
│   ├── check-balances.sh    # Verificar saldos de todos os nós
│   ├── check-lastblock.sh   # Verificar altura da blockchain de todos os nós
│   ├── network-status.sh    # Dashboard de status da rede
│   ├── docker-test.sh       # Teste automatizado da rede Docker
│   └── demo.sh              # Script de demonstração rápida
├── docker-compose.yml       # Setup de rede multi-nó Docker
├── Dockerfile               # Definição de imagem do container
├── go.mod                   # Dependências de módulos Go
├── go.sum                   # Checksums de módulos Go
├── Makefile                 # Automação de build
├── LICENSE                  # Licença MIT
├── README.md                # README em inglês
└── README.pt-br.md         # Este arquivo
```

## 🚀 Começando

### Pré-requisitos

- Go 1.22 ou superior
- Docker & Docker Compose (para testes multi-nó)

### Quick Start (Docker - Recomendado)

A rede Docker já vem pré-configurada com 4 nós e é a forma mais fácil de testar:

```bash
# Clone o repositório
git clone https://github.com/marcocsrachid/blockchain-go.git
cd blockchain-go

# Build e inicie a rede (4 nós: 1 seed, 2 miners, 1 regular)
docker-compose build
docker-compose up -d

# Verifique o status
docker-compose ps

# Veja os logs
docker-compose logs -f
```

**Portas expostas:**
- `4000` - Seed Node (API HTTP)
- `4001` - Miner 1 (API HTTP)
- `4002` - Miner 2 (API HTTP)
- `4003` - Regular Node (API HTTP)

**Scripts úteis:**
```bash
# Status completo da rede
./scripts/network-status.sh

# Verificar altura dos blocos
./scripts/check-lastblock.sh

# Verificar saldos
./scripts/check-balances.sh
```

### Build Manual (Local)

Se você quiser compilar e executar localmente:

#### 1. Build do binário

```bash
# Build padrão
go build -o build/blockchain cmd/blockchain/main.go

# Ou build estático (para Docker Alpine)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix netgo -ldflags '-s -w' -o build/blockchain cmd/blockchain/main.go
```

#### 2. Criar wallet e blockchain

```bash
# Criar uma wallet (anote o endereço gerado)
./build/blockchain createwallet

# Listar endereços
./build/blockchain listaddresses

# Criar a blockchain com endereço de recompensa
./build/blockchain createblockchain -address SEU_ENDERECO
```

#### 3. Startar um node

**Node minerador (produz blocos):**
```bash
# Terminal 1 - Seed/Miner Node
NODE_ID=node1 ./build/blockchain startnode -port 3000 -miner SEU_ENDERECO
```

**Node regular (não minera):**
```bash
# Terminal 2 - Regular Node (conecta ao node1)
NODE_ID=node2 SEED_NODE=localhost:3000 ./build/blockchain startnode -port 3001
```

**Variáveis de ambiente importantes:**
- `NODE_ID` - ID único do node (define o diretório de dados)
- `SEED_NODE` - Endereço do seed node para conectar
- `-port` - Porta P2P do node (default: 3000)
- `-apiport` - Porta da API HTTP (default: 4000)
- `-miner` - Endereço para receber recompensas (ativa mineração)

### Usando a API HTTP

Todos os nodes expõem uma API REST:

```bash
# Verificar status da rede
curl http://localhost:4000/api/networkinfo | jq

# Listar endereços
curl http://localhost:4000/api/addresses | jq

# Verificar saldo
curl http://localhost:4000/api/balance/SEU_ENDERECO | jq

# Criar nova wallet
curl -X POST http://localhost:4000/api/createwallet | jq

# Enviar transação
curl -X POST http://localhost:4000/api/send \
  -H "Content-Type: application/json" \
  -d '{
    "from": "ENDERECO_ORIGEM",
    "to": "ENDERECO_DESTINO",
    "amount": 10
  }' | jq

# Ver último bloco
curl http://localhost:4000/api/lastblock | jq

# Listar peers conhecidos
curl http://localhost:4000/api/peers | jq
```

### Exemplo Completo (3 Nodes)

```bash
# Terminal 1 - Seed Node (não minera, apenas coordena)
NODE_ID=seed ./build/blockchain createblockchain -address 1SeedAddress...
NODE_ID=seed ./build/blockchain startnode -port 3000

# Terminal 2 - Miner 1
NODE_ID=miner1 SEED_NODE=localhost:3000 ./build/blockchain startnode -port 3001 -apiport 4001 -miner 1Miner1Address...

# Terminal 3 - Miner 2
NODE_ID=miner2 SEED_NODE=localhost:3000 ./build/blockchain startnode -port 3002 -apiport 4002 -miner 1Miner2Address...

# Terminal 4 - Enviar transação via API
curl -X POST http://localhost:4001/api/send \
  -H "Content-Type: application/json" \
  -d '{"from":"1Miner1Address...","to":"1Miner2Address...","amount":50}' | jq

# Aguarde ~60-90s para mineração...

# Verificar saldos
curl http://localhost:4001/api/balance/1Miner1Address... | jq
curl http://localhost:4002/api/balance/1Miner2Address... | jq
```

### Acessando Containers Docker

```bash
# Executar comandos dentro dos containers
docker exec -it blockchain-seed /app/blockchain listaddresses
docker exec -it blockchain-miner1 /app/blockchain listaddresses

# Ver logs de um node específico
docker-compose logs -f node-seed
docker-compose logs -f node-miner1

# Parar a rede
docker-compose down

# Parar e limpar dados (reset completo)
docker-compose down -v
```

📖 Para detalhes completos sobre a implementação de rede, veja [docs/NETWORK.md](docs/NETWORK.md) e [docs/NETWORK.pt-br.md](docs/NETWORK.pt-br.md)

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

| Funcionalidade | Bitcoin | Este Projeto | Status |
|---------------|---------|--------------|--------|
| Proof of Work | ✅ | ✅ | Implementado |
| Merkle Tree | ✅ | ✅ | Implementado |
| UTXO Model | ✅ | ✅ | Implementado |
| ECDSA | secp256k1 | P256 | Implementado (curva diferente) |
| Base58 | ✅ | ✅ | Implementado |
| Endereços | ✅ | ✅ | Implementado |
| Transações | ✅ | ✅ | Implementado (simplificado) |
| Rede P2P | ✅ Completa | ✅ Básica | Implementado (sem DNS seeds) |
| Mempool | ✅ | ✅ | Implementado (sem RBF) |
| HTTP API | ❌ | ✅ | Extra: REST API |
| Scripts | ✅ | ❌ | Não implementado |
| Ajuste Dificuldade | ✅ A cada 2016 blocos | ❌ Fixa | Simplificado |
| Halving | ✅ | ✅ | Implementado |

**Similaridade com Bitcoin: ~93%** dos conceitos fundamentais

### Diferenças Principais

1. **Dificuldade Fixa**: Não ajusta automaticamente a cada 2016 blocos
2. **Halving Simplificado**: Implementado mas sem complexidade de ajuste
3. **P2P Básico**: Sem DNS seeds, descoberta manual de peers
4. **Sem Scripts**: Não usa linguagem Script para condições de gasto
5. **Mempool Básico**: Sem priority fees ou Replace-By-Fee (RBF)
6. **Sem SPV**: Simplified Payment Verification não implementado
7. **Sem SegWit**: Segregated Witness não implementado
8. **Sem Lightning**: Lightning Network não implementado
9. **API REST**: Extra não presente no Bitcoin Core (tem RPC)

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

