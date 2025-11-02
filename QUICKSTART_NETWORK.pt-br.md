# Guia Rápido de Rede

> **English version**: [QUICKSTART_NETWORK.md](QUICKSTART_NETWORK.md)

## 🚀 Teste Rápido com Docker (Recomendado)

### 1. Iniciar a Rede

```bash
# Limpar dados antigos e compilar
make docker-clean
make docker-build

# Iniciar rede com 4 nós
make docker-up
```

Isso inicia:
- **Nó Seed** (localhost:3000) - Não minera
- **Minerador 1** (localhost:3001) - Nó minerador
- **Minerador 2** (localhost:3002) - Nó minerador  
- **Nó Regular** (localhost:3003) - Não minera

### 2. Acompanhar a Rede

```bash
# Ver todos os logs
make docker-logs

# Ver nó específico
docker-compose logs -f node-miner1
```

### 3. Interagir com a Rede

```bash
# Listar carteiras no minerador 1
docker exec -it blockchain-miner1 /app/blockchain listaddresses

# Obter um endereço
ADDR=$(docker exec -it blockchain-miner1 /app/blockchain listaddresses | head -1 | tr -d '\r')

# Verificar saldo
docker exec -it blockchain-miner1 /app/blockchain getbalance -address "$ADDR"

# Ver blockchain
docker exec -it blockchain-seed /app/blockchain printchain
```

### 4. Parar a Rede

```bash
# Parar containers
make docker-down

# Parar e remover todos os dados
make docker-clean
```

## 🖥️ Teste Manual (Múltiplos Terminais)

### Terminal 1: Nó Seed

```bash
# Limpar e compilar
make clean
make build

# Criar blockchain
./build/blockchain createwallet
# Salve o endereço: ADDRESS1=<seu_endereço>

./build/blockchain createblockchain -address <ADDRESS1>

# Iniciar nó seed
./build/blockchain startnode -port 3000 -miner <ADDRESS1>
```

### Terminal 2: Nó Minerador

```bash
# Criar carteira
./build/blockchain createwallet
# Salve o endereço: ADDRESS2=<seu_endereço>

# Iniciar nó minerador
./build/blockchain startnode -port 3001 -miner <ADDRESS2>
```

### Terminal 3: Nó Regular

```bash
# Iniciar nó regular
./build/blockchain startnode -port 3002
```

### Terminal 4: Enviar Transações

```bash
# Enviar transação
./build/blockchain send -from <ADDRESS1> -to <ADDRESS2> -amount 10

# Verificar saldos
./build/blockchain getbalance -address <ADDRESS1>
./build/blockchain getbalance -address <ADDRESS2>

# Ver blockchain
./build/blockchain printchain
```

## 📊 O Que Esperar

1. **Inicialização do Nó**
   - Nós conectam ao nó seed (localhost:3000)
   - Troca de versão e altura da blockchain
   - Sincronização da blockchain

2. **Fluxo de Transação**
   - Transação criada em qualquer nó
   - Broadcast para todos os peers
   - Adicionada ao mempool

3. **Mineração**
   - Mineradores coletam transações do mempool
   - Mineram novo bloco com PoW
   - Broadcast do novo bloco para a rede
   - Todos os nós validam e adicionam o bloco

4. **Consenso**
   - Todos os nós mantêm a mesma blockchain
   - Regra da cadeia mais longa (como Bitcoin)
   - Resolução automática de forks

## 🔍 Dicas de Depuração

### Verificar Conectividade do Nó

```bash
# Testar se seed está escutando
nc -zv localhost 3000

# Ver peers conectados
./build/blockchain peers
```

### Ver Logs do Docker

```bash
# Todos os nós
docker-compose logs

# Intervalo de tempo específico (últimos 10 minutos)
docker-compose logs --since 10m

# Acompanhar ao vivo
docker-compose logs -f node-miner1
```

### Acessar Container Docker

```bash
# Shell interativo
docker exec -it blockchain-seed sh

# Executar comandos
docker exec -it blockchain-seed /app/blockchain printchain
```

## 🎯 Cenários de Teste

### Cenário 1: Rede Básica

1. Iniciar nó seed + 2 mineradores
2. Criar transação
3. Observar blocos sendo minerados
4. Verificar se todos os nós têm a mesma blockchain

### Cenário 2: Nó Atrasado

1. Iniciar seed + minerador 1
2. Minerar vários blocos
3. Iniciar minerador 2 (nó atrasado)
4. Verificar se minerador 2 sincroniza a blockchain

### Cenário 3: Múltiplas Transações

1. Criar 3 carteiras
2. Enviar múltiplas transações
3. Observar mineradores competindo
4. Verificar se todos os UTXOs estão corretos

## 📝 Problemas Comuns

### Porta Já em Uso

```bash
# Encontrar processo
lsof -i :3000

# Matar processo
kill -9 <PID>

# Ou usar porta diferente
./build/blockchain startnode -port 3005
```

### Blockchain Não Sincroniza

```bash
# Reindexar UTXO
./build/blockchain reindexutxo

# Ou limpar e reiniciar
rm -rf ./tmp
./build/blockchain createblockchain -address <ENDEREÇO>
```

### Build Docker Falha

```bash
# Limpar cache do Docker
docker system prune -a

# Recompilar
make docker-build
```

## 🎓 Exercícios de Aprendizado

1. **Modificar Dificuldade de Mineração**
   - Editar `internal/blockchain/proof.go`
   - Mudar constante `Difficulty`
   - Observar mudanças no tempo de mineração

2. **Mudar Threshold do Mempool**
   - Editar `internal/network/server.go`
   - Modificar `len(memoryPool) >= 2`
   - Testar batching diferente de transações

3. **Adicionar Estatísticas de Rede**
   - Rastrear mensagens enviadas/recebidas
   - Monitorar conexões de peers
   - Registrar tempo de sincronização

4. **Implementar Persistência**
   - Salvar lista de peers em disco
   - Restaurar conexões ao reiniciar
   - Adicionar pontuação de reputação de peers

## 📚 Próximos Passos

- Ler [docs/NETWORK.pt-br.md](docs/NETWORK.pt-br.md) para documentação detalhada
- Explorar diferenças do protocolo P2P do Bitcoin
- Implementar funcionalidades adicionais:
  - Relay de blocos compactos
  - Prioridade de transações no mempool
  - Protocolo de descoberta de peers
  - Dashboard de estatísticas de rede

## 🐛 Solução de Problemas

Se algo não funcionar:

1. Verificar se blockchain existe: `ls -la tmp/`
2. Verificar se portas estão disponíveis: `netstat -an | grep 3000`
3. Ver logs: `make docker-logs`
4. Limpar tudo: `make docker-clean && rm -rf tmp/`
5. Começar do zero: Seguir "Teste Rápido com Docker" desde o passo 1

---

**Feliz Networking Blockchain! 🎉**

