# Documentação da Rede Blockchain

## Visão Geral

Esta implementação de blockchain inclui uma camada de rede peer-to-peer (P2P) que permite que múltiplos nós se comuniquem, sincronizem e mantenham um ledger distribuído.

## Arquitetura da Rede

### Componentes

1. **Servidor** (`internal/network/server.go`)
   - Servidor TCP para gerenciar conexões
   - Gerencia conexões de peers e roteamento de mensagens
   - Coordena sincronização da blockchain

2. **Gerenciamento de Peers** (`internal/network/peer.go`)
   - Mantém lista de peers conectados
   - Gerencia ciclo de vida das conexões
   - Operações thread-safe

3. **Protocolo** (`internal/network/protocol.go`)
   - Define tipos e formatos de mensagens
   - Serialização/deserialização usando gob
   - Roteamento de comandos

### Tipos de Mensagem

| Comando | Descrição |
|---------|-----------|
| `version` | Handshake inicial entre peers |
| `getblocks` | Solicita hashes da blockchain |
| `inv` | Lista blocos/transações disponíveis |
| `getdata` | Solicita bloco ou transação específica |
| `block` | Envia dados de bloco |
| `tx` | Envia dados de transação |
| `addr` | Compartilha endereços de peers |
| `ping/pong` | Mensagens keep-alive |

## Fluxo da Rede

### 1. Inicialização do Nó

```
Nó A inicia → Escuta na porta → Conecta aos seed nodes → Troca version
```

### 2. Sincronização da Blockchain

```
Nó A ← version (altura: 10) ← Nó B (seed)
Nó A → getblocks → Nó B
Nó A ← inv (hashes de blocos) ← Nó B
Nó A → getdata (hash do bloco) → Nó B
Nó A ← block (dados) ← Nó B
```

### 3. Broadcast de Transações

```
Nó A cria transação → broadcast para todos peers
Nó B recebe transação → adiciona ao mempool
Nó B (minerador) minera bloco → broadcast novo bloco
Nó A recebe novo bloco → valida → adiciona à chain
```

## Executando a Rede

### Teste Local

#### Método 1: Múltiplos Terminais

**Terminal 1 - Nó Seed + Minerador:**
```bash
./build/blockchain startnode -port 3000 -miner <ENDEREÇO>
```

**Terminal 2 - Nó Minerador:**
```bash
./build/blockchain startnode -port 3001 -miner <ENDEREÇO>
```

**Terminal 3 - Nó Regular:**
```bash
./build/blockchain startnode -port 3002
```

### Teste com Docker

#### Início Rápido

```bash
make docker-build    # Constrói imagens
make docker-up       # Inicia rede (4 nós)
make docker-logs     # Visualiza logs
```

#### Comandos Docker Manuais

```bash
# Iniciar rede
docker-compose up -d

# Ver logs
docker-compose logs -f node-seed

# Executar comandos nos containers
docker exec -it blockchain-seed /app/blockchain listaddresses
docker exec -it blockchain-miner1 /app/blockchain getbalance -address <ENDEREÇO>

# Parar rede
docker-compose down

# Limpar (remove dados)
docker-compose down -v
```

## Configuração da Rede

### Rede Docker Compose

O `docker-compose.yml` cria uma rede com:

- **Nó Seed** (172.20.0.2:3000) - Não minera, nó seed
- **Minerador 1** (172.20.0.3:3001) - Nó minerador
- **Minerador 2** (172.20.0.4:3002) - Nó minerador
- **Nó Regular** (172.20.0.5:3003) - Não minera

## Comandos CLI

### Iniciar Nó

```bash
# Nó minerador
./build/blockchain startnode -port 3000 -miner <ENDEREÇO_CARTEIRA>

# Nó regular
./build/blockchain startnode -port 3000
```

### Adicionar Peer

```bash
./build/blockchain addpeer -address localhost:3001
```

### Listar Peers

```bash
./build/blockchain peers
```

## Mineração em Modo Rede

Quando um nó é iniciado com a flag `-miner`:

1. Nó recebe transações via P2P
2. Transações são adicionadas ao mempool
3. Quando mempool atinge threshold (2+ transações), mineração inicia
4. Novo bloco é minerado e broadcast para todos peers
5. Peers validam e adicionam bloco às suas chains

## Resolução de Problemas

### Nó Não Conecta

**Problema:** Nó não consegue conectar ao seed

**Solução:**
```bash
# Verificar se seed está rodando
nc -zv localhost 3000

# Verificar logs
docker logs blockchain-seed
```

### Blockchain Dessincronizada

**Problema:** Nó tem blockchain diferente dos peers

**Solução:**
```bash
# Reindexar UTXO
./build/blockchain reindexutxo

# Ou reiniciar nó
docker-compose restart node-miner1
```

## Considerações de Segurança

⚠️ **Esta é uma implementação educacional. Para produção:**

1. **Adicionar TLS/SSL** - Criptografar conexões
2. **Autenticação** - Verificar identidade dos peers
3. **Proteção DDoS** - Rate limiting, limite de conexões
4. **Prevenção Sybil Attack** - PoW para descoberta de peers
5. **Prevenção Eclipse Attack** - Seleção diversa de peers
6. **Validação de Input** - Sanitizar mensagens de rede

## Comparação com Protocolo Bitcoin

| Feature | Bitcoin | Esta Implementação |
|---------|---------|-------------------|
| Protocolo | Binário customizado | Gob encoding |
| Descoberta | DNS seeds | Nó seed estático |
| Conexão | Persistente | Por mensagem |
| Mempool | Sim | Sim (em memória) |
| Block relay | Compact blocks | Blocos completos |
| SPV | Sim | Não |

## Próximos Passos

Para tornar esta rede mais robusta:

1. **Conexões Persistentes** - Manter conexões abertas
2. **Descoberta de Peers** - Encontrar peers automaticamente
3. **Headers First** - Sincronização mais rápida
4. **Suporte SPV** - Clientes leves
5. **Estatísticas de Rede** - Monitorar saúde dos peers
6. **Reconexão Automática** - Gerenciar falhas de rede

---

**Veja também:**
- [QUICKSTART_NETWORK.md](../QUICKSTART_NETWORK.md) - Guia rápido em português
- [NETWORK.md](NETWORK.md) - Documentação completa em inglês

