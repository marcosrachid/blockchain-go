# Mineração e Ajuste de Dificuldade

## Implementação Atual

### Proof of Work (PoW)
- **Dificuldade Fixa**: 18 (definida em `config.go`)
- Mineradores computam hashes SHA-256 continuamente até encontrar um hash com 18 zeros iniciais
- Sem timers - mineração acontece **continuamente** até um bloco válido ser encontrado

### Como a Mineração Funciona

1. **Mineração Contínua**: Quando um nó é configurado como minerador, inicia um loop infinito
2. **Criação do Bloco**: Coleta transações do mempool + transação coinbase
3. **Computação do PoW**: Tenta diferentes valores de nonce até o hash atender à dificuldade
4. **Broadcast**: Uma vez encontrado, transmite imediatamente para todos os peers
5. **Repetir**: Inicia imediatamente a mineração do próximo bloco

```go
// Loop de mineração (simplificado)
for {
    // Coletar transações
    txs := getTransactionsFromMempool()
    coinbase := createCoinbaseTransaction()
    
    // Minerar bloco (PoW)
    block := mineBlock(txs + coinbase) // Leva tempo baseado na dificuldade
    
    // Broadcast
    broadcastBlock(block)
}
```

## Ajuste de Dificuldade do Bitcoin

No Bitcoin, a dificuldade ajusta **automaticamente** a cada 2016 blocos (~2 semanas):

```
nova_dificuldade = dificuldade_atual * (tempo_real / tempo_alvo)
```

- **Alvo**: 10 minutos por bloco
- **Intervalo de Ajuste**: 2016 blocos
- Se blocos foram mais rápidos → aumenta dificuldade
- Se blocos foram mais lentos → diminui dificuldade

### Exemplo

Se 2016 blocos levaram **1 semana** em vez de **2 semanas**:
- Blocos foram minerados **2x mais rápido** que o alvo
- Nova dificuldade = atual × (1 semana / 2 semanas) = atual × 0.5
- **Dificuldade aumenta** para desacelerar a mineração

## Melhorias Futuras

### 1. Ajuste Dinâmico de Dificuldade

Adicionar em `internal/blockchain/config.go`:

```go
const (
    DifficultyAdjustmentInterval = 100 // Ajustar a cada 100 blocos
    TargetBlockTime = 60 // segundos
)
```

Implementar em `blockchain.go`:

```go
func (chain *Blockchain) CalculateNewDifficulty() int {
    lastBlock := chain.GetLastBlock()
    
    // Pegar bloco do intervalo de ajuste atrás
    targetBlock := chain.GetBlockAtHeight(lastBlock.Height - DifficultyAdjustmentInterval)
    
    // Calcular tempo real
    actualTime := lastBlock.Timestamp - targetBlock.Timestamp
    expectedTime := DifficultyAdjustmentInterval * TargetBlockTime
    
    // Ajustar dificuldade
    if actualTime < expectedTime / 2 {
        // Blocos muito rápidos, aumentar dificuldade
        return currentDifficulty + 1
    } else if actualTime > expectedTime * 2 {
        // Blocos muito lentos, diminuir dificuldade
        return max(currentDifficulty - 1, 1)
    }
    
    return currentDifficulty
}
```

### 2. Competição de Mineração

Comportamento atual com múltiplos mineradores:
- Todos os mineradores competem para encontrar o próximo bloco
- Primeiro a encontrar PoW válido transmite e vence
- Outros descartam seu trabalho e começam no próximo bloco
- Este é o comportamento **correto** similar ao Bitcoin

### 3. Blocos Órfãos

Quando dois mineradores encontram blocos simultaneamente:
- Ambos transmitem seus blocos
- Rede pode temporariamente se dividir
- Próximo bloco determina o vencedor (regra da cadeia mais longa)
- Bloco perdedor torna-se "órfão" e é descartado

## Por Que Não Usar Timer?

❌ **Abordagem Errada** (implementação anterior):
```go
ticker := time.NewTicker(60 * time.Second)
for {
    <-ticker.C
    mineBlock() // Força mineração a cada 60 segundos
}
```

✅ **Abordagem Correta** (atual):
```go
for {
    block := mineBlock() // Leva tempo naturalmente baseado na dificuldade
    broadcast(block)
    // Inicia imediatamente próximo bloco
}
```

A **dificuldade** controla o tempo, não um timer!

## Configuração

Para ajustar o tempo alvo entre blocos, modifique em `config.go`:

```go
const (
    Difficulty = 18  // Maior = blocos mais lentos, mais segurança
                     // Menor = blocos mais rápidos, menos segurança
)
```

- Dificuldade 12: ~0.1 segundos por bloco
- Dificuldade 18: ~1 minuto por bloco (atual)
- Dificuldade 20: ~4 minutos por bloco
- Dificuldade 24: ~1 hora por bloco

## Referências

- [Bitcoin Difficulty Adjustment](https://en.bitcoin.it/wiki/Difficulty)
- [Proof of Work](https://en.bitcoin.it/wiki/Proof_of_work)
- [Mining](https://en.bitcoin.it/wiki/Mining)

