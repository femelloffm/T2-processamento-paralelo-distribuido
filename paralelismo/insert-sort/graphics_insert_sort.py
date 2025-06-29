import matplotlib.pyplot as plt
import numpy as np

# Dados de tempo de execução (em segundos)
dados = {
    300: {1: 0.003807, 2: 0.003014, 3: 0.003087, 4: 0.002501},
    600: {1: 0.015935, 2: 0.012831, 3: 0.009037, 4: 0.008715},
    900: {1: 0.035413, 2: 0.025181, 3: 0.019973, 4: 0.017854}
}

# Número de processadores
processadores = [1, 2, 3, 4]

# Configuração da figura com 3 subplots
fig, axes = plt.subplots(1, 3, figsize=(15, 5))
fig.suptitle('Speedup vs Número de Processadores', fontsize=16, fontweight='bold')

# Cores para cada tamanho
cores = ['#2E8B57', '#4169E1', '#DC143C']  # Verde, Azul, Vermelho
tamanhos = [300, 600, 900]

for i, n in enumerate(tamanhos):
    # Cálculo do speedup para cada número de processadores
    tempos = [dados[n][p] for p in processadores]
    tempo_sequencial = dados[n][1]  # Tempo com 1 processador como referência
    speedups = [tempo_sequencial / tempo for tempo in tempos]
    
    # Criar o gráfico
    ax = axes[i]
    
    # Plot dos dados
    ax.plot(processadores, speedups, 'o-', color=cores[i], linewidth=2.5, 
            markersize=8, markerfacecolor=cores[i], markeredgecolor='white', 
            markeredgewidth=2)
    
    # Linha de speedup ideal (linear)
    speedup_ideal = processadores
    ax.plot(processadores, speedup_ideal, '--', color='gray', alpha=0.7, 
            linewidth=1.5, label='Speedup Ideal')
    
    # Configurações do gráfico
    ax.set_title(f'{n} números', fontsize=14, fontweight='bold', pad=15)
    ax.set_xlabel('Número de Processadores', fontsize=12)
    ax.set_ylabel('Speedup', fontsize=12)
    ax.grid(True, alpha=0.3, linestyle='-', linewidth=0.5)
    ax.set_xlim(0.8, 4.2)
    ax.set_ylim(0.8, max(max(speedups), 4) + 0.2)
    
    # Configurar ticks
    ax.set_xticks(processadores)
    ax.set_xticklabels(processadores)
    
    # Adicionar valores nos pontos
    for j, (p, s) in enumerate(zip(processadores, speedups)):
        ax.annotate(f'{s:.2f}', (p, s), textcoords="offset points", 
                   xytext=(0,10), ha='center', fontsize=10, fontweight='bold')
    
    # Adicionar legenda apenas no primeiro gráfico
    if i == 0:
        ax.legend(loc='upper left', fontsize=10)

# Ajustar layout
plt.tight_layout()

# Mostrar estatísticas
print("=== ANÁLISE DE SPEEDUP ===\n")
for n in tamanhos:
    print(f"Problema com {n} números:")
    tempo_base = dados[n][1]
    print(f"  Tempo base (1 proc): {tempo_base:.6f}s")
    
    for p in processadores[1:]:  # Pular 1 processador
        tempo_atual = dados[n][p]
        speedup = tempo_base / tempo_atual
        eficiencia = speedup / p * 100
        print(f"  {p} proc: {tempo_atual:.6f}s | Speedup: {speedup:.2f} | Eficiência: {eficiencia:.1f}%")
    print()

# Calcular e mostrar speedup médio
print("=== SPEEDUP MÉDIO POR NÚMERO DE PROCESSADORES ===")
for p in processadores[1:]:
    speedups_medios = []
    for n in tamanhos:
        speedup = dados[n][1] / dados[n][p]
        speedups_medios.append(speedup)
    speedup_medio = np.mean(speedups_medios)
    print(f"{p} processadores: {speedup_medio:.2f}")

plt.show()