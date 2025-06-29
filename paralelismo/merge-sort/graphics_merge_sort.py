import matplotlib.pyplot as plt
import numpy as np

dados = {
    1: {
        50000: {1: 0.116638, 2: 0.068143, 4: 0.045884},
        100000: {1: 0.215335, 2: 0.164258, 4: 0.085965},
        500000: {1: 1.239434, 2: 0.605324, 4: 0.376127},
        1000000: {1: 2.395317, 2: 1.319387, 4: 0.847059}
    },
    500: {
        50000: {1: 0.008533, 2: 0.005924, 4: 0.003907},
        100000: {1: 0.017033, 2: 0.009607, 4: 0.007124},
        500000: {1: 0.086655, 2: 0.058231, 4: 0.034174},
        1000000: {1: 0.195443, 2: 0.125859, 4: 0.066597}
    },
    1000: {
        50000: {1: 0.008000, 2: 0.004514, 4: 0.003774},
        100000: {1: 0.020448, 2: 0.010387, 4: 0.008529},
        500000: {1: 0.094249, 2: 0.049631, 4: 0.033649},
        1000000: {1: 0.190690, 2: 0.092126, 4: 0.054412}
    },
    5000: {
        50000: {1: 0.008762, 2: 0.005330, 4: 0.004470},
        100000: {1: 0.019760, 2: 0.010317, 4: 0.006780},
        500000: {1: 0.094068, 2: 0.055543, 4: 0.028345},
        1000000: {1: 0.181873, 2: 0.107390, 4: 0.068355}
    }
}

# Parametros
processadores = [1, 2, 4]
granularidade = [1, 500, 1000, 5000]
tamanhos = [50000, 100000, 500000, 1000000]

# Configuração da figura com 4 subplots
fig, axes = plt.subplots(len(granularidade), len(tamanhos), figsize=(18, 16))
fig.suptitle('Speedup vs Número de Processadores', fontsize=16, fontweight='bold')

# Cores para cada tamanho
cores = ['#2E8B57', '#4169E1', "#DC143C", '#7814dc']  # Verde, Azul, Vermelho, Roxo

for indexG, g in enumerate(granularidade):
    dados_granularidade = dados[g]
    for i, n in enumerate(tamanhos):
        # Cálculo do speedup para cada número de processadores
        tempos = [dados_granularidade[n][p] for p in processadores]
        tempo_sequencial = dados_granularidade[n][1]  # Tempo com 1 processador como referência
        speedups = [tempo_sequencial / tempo for tempo in tempos]
        
        # Criar o gráfico
        ax = axes[indexG, i]
        
        # Plot dos dados
        ax.plot(processadores, speedups, 'o-', color=cores[i], linewidth=2.5, 
                markersize=8, markerfacecolor=cores[i], markeredgecolor='white', 
                markeredgewidth=2)
        
        # Linha de speedup ideal (linear)
        speedup_ideal = processadores
        ax.plot(processadores, speedup_ideal, '--', color='gray', alpha=0.7, 
                linewidth=1.5, label='Speedup Ideal')
        
        # Configurações do gráfico
        ax.set_title(f'Granularidade {g} - Tamanho {n}', fontsize=14, fontweight='bold', pad=15)
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

    # Mostrar estatísticas
    print("=== ANÁLISE DE SPEEDUP ===\n")
    for n in tamanhos:
        print(f"Problema com {n} números:")
        tempo_base = dados_granularidade[n][1]
        print(f"  Tempo base (1 proc): {tempo_base:.6f}s")
        
        for p in processadores[1:]:  # Pular 1 processador
            tempo_atual = dados_granularidade[n][p]
            speedup = tempo_base / tempo_atual
            eficiencia = speedup / p * 100
            print(f"  {p} proc: {tempo_atual:.6f}s | Speedup: {speedup:.2f} | Eficiência: {eficiencia:.1f}%")
        print()

    # Calcular e mostrar speedup médio
    print("=== SPEEDUP MÉDIO POR NÚMERO DE PROCESSADORES ===")
    for p in processadores[1:]:
        speedups_medios = []
        for n in tamanhos:
            speedup = dados_granularidade[n][1] / dados_granularidade[n][p]
            speedups_medios.append(speedup)
        speedup_medio = np.mean(speedups_medios)
        print(f"{p} processadores: {speedup_medio:.2f}")


# Ajustar layout
plt.tight_layout()
plt.show()