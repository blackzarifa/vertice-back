-- Create database is handled by the Go application
USE bancovertice_db;

-- Endereco must come first since usuario references it
CREATE TABLE IF NOT EXISTS endereco(
	id_endereco INT PRIMARY KEY AUTO_INCREMENT,
    cep VARCHAR(10) NOT NULL,
    local VARCHAR(100) NOT NULL,
    numero_casa INT NOT NULL,
    bairro VARCHAR(50) NOT NULL,
    cidade VARCHAR(50) NOT NULL,
    estado CHAR(2) NOT NULL,
    complemento VARCHAR(50),
    INDEX (cep)
);

CREATE TABLE IF NOT EXISTS usuario(
	id_usuario INT PRIMARY KEY AUTO_INCREMENT,
    id_endereco INT,
    nome VARCHAR(100) NOT NULL,
    cpf VARCHAR(11) UNIQUE NOT NULL,
    data_nascimento DATE NOT NULL,
    telefone VARCHAR(15) NOT NULL,
    tipo_usuario ENUM('FUNCIONARIO', 'CLIENTE') NOT NULL,
    senha_hash VARCHAR(32) NOT NULL,
    otp_ativo VARCHAR(6),
    otp_expiracao DATETIME,
    FOREIGN KEY (id_endereco) REFERENCES endereco(id_endereco)
);

CREATE TABLE IF NOT EXISTS funcionario(
    id_funcionario INT PRIMARY KEY AUTO_INCREMENT,
    id_usuario INT,
    id_supervisor INT,
    codigo_funcionario VARCHAR(20) UNIQUE NOT NULL,
    cargo ENUM('ESTAGIARIO', 'ATENDENTE', 'GERENTE') NOT NULL,
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario),
    FOREIGN KEY (id_supervisor) REFERENCES funcionario(id_funcionario)
);

CREATE TABLE IF NOT EXISTS cliente(
	id_cliente INT PRIMARY KEY AUTO_INCREMENT,
    id_usuario INT,
    score_credito DECIMAL(5,2) DEFAULT 0,
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario)
);

CREATE TABLE IF NOT EXISTS agencia(
	id_agencia INT PRIMARY KEY AUTO_INCREMENT,
    nome VARCHAR(50) NOT NULL,
    codigo_agencia VARCHAR(10) UNIQUE NOT NULL,
	endereco_id INT,
    FOREIGN KEY (endereco_id) REFERENCES endereco(id_endereco)
);

CREATE TABLE IF NOT EXISTS conta(
    id_conta INT PRIMARY KEY AUTO_INCREMENT,
    numero_conta VARCHAR(20) UNIQUE,
    digito_verificador CHAR(1),
    id_agencia INT,
    saldo DECIMAL(15,2) NOT NULL DEFAULT 0,
    tipo_conta ENUM('POUPANCA', 'CORRENTE', 'INVESTIMENTO') NOT NULL,
    id_cliente INT,
    data_abertura DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status ENUM('ATIVA', 'ENCERRADA', 'BLOQUEADA') NOT NULL DEFAULT 'ATIVA',
    FOREIGN KEY (id_agencia) REFERENCES agencia(id_agencia),
    FOREIGN KEY (id_cliente) REFERENCES cliente(id_cliente),
    INDEX (numero_conta)
);

CREATE TABLE IF NOT EXISTS conta_poupanca(
    id_conta_poupanca INT AUTO_INCREMENT PRIMARY KEY,
    id_conta INT UNIQUE,
    taxa_rendimento DECIMAL(5,2) NOT NULL,
    ultimo_rendimento DATETIME,
    FOREIGN KEY (id_conta) REFERENCES conta(id_conta)
);

CREATE TABLE IF NOT EXISTS conta_corrente(
    id_conta_corrente INT AUTO_INCREMENT PRIMARY KEY,
    id_conta INT UNIQUE,
    limite DECIMAL(15,2) NOT NULL DEFAULT 0,
    data_vencimento DATE NOT NULL,
    taxa_manutencao DECIMAL(5,2) NOT NULL DEFAULT 0,
    FOREIGN KEY (id_conta) REFERENCES conta(id_conta)
);

CREATE TABLE IF NOT EXISTS conta_investimento(
    id_conta_investimento INT PRIMARY KEY AUTO_INCREMENT,
    id_conta INT UNIQUE,
    perfil_risco ENUM('BAIXO', 'MEDIO', 'ALTO') NOT NULL,
    valor_minimo DECIMAL(15,2) NOT NULL,
    taxa_rendimento_base DECIMAL(5,2) NOT NULL,
    FOREIGN KEY (id_conta) REFERENCES conta(id_conta)
);

CREATE TABLE IF NOT EXISTS transacao(
    id_transacao INT PRIMARY KEY AUTO_INCREMENT,
    id_conta_origem INT,
    id_conta_destino INT,
    tipo_transacao ENUM('DEPOSITO', 'SAQUE', 'TRANSFERENCIA', 'TAXA', 'RENDIMENTO') NOT NULL,
    valor DECIMAL(15,2) NOT NULL,
    data_hora TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    descricao VARCHAR(100),
    FOREIGN KEY (id_conta_origem) REFERENCES conta(id_conta),
    FOREIGN KEY (id_conta_destino) REFERENCES conta(id_conta),
    INDEX (data_hora)
);

CREATE TABLE IF NOT EXISTS auditoria(
    id_auditoria INT PRIMARY KEY AUTO_INCREMENT,
    id_usuario INT,
    acao VARCHAR(50) NOT NULL,
    data_hora TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    detalhes TEXT,
    FOREIGN KEY (id_usuario) REFERENCES usuario(id_usuario)
);

CREATE TABLE IF NOT EXISTS relatorio(
    id_relatorio INT PRIMARY KEY AUTO_INCREMENT,
    id_funcionario INT,
    tipo_relatorio VARCHAR(50) NOT NULL,
    data_geracao TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    conteudo TEXT NOT NULL,
    FOREIGN KEY (id_funcionario) REFERENCES funcionario(id_funcionario)
);
