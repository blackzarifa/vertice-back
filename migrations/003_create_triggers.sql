-- Triggers

DELIMITER //

DROP TRIGGER IF EXISTS gerar_numero_conta//
CREATE TRIGGER gerar_numero_conta
AFTER INSERT ON conta
FOR EACH ROW
BEGIN
    DECLARE novo_numero VARCHAR(20);
    SET novo_numero = LPAD(NEW.id_conta, 10, '0');
    UPDATE conta 
    SET numero_conta = novo_numero,
        digito_verificador = calcular_dv(novo_numero)
    WHERE id_conta = NEW.id_conta;
END//

DROP TRIGGER IF EXISTS registrar_auditoria//
CREATE TRIGGER registrar_auditoria
AFTER INSERT ON conta
FOR EACH ROW
BEGIN
    DECLARE usuario_dono INT;
    SELECT id_usuario INTO usuario_dono FROM cliente WHERE id_cliente = NEW.id_cliente;
    
    INSERT INTO auditoria (id_usuario, acao, detalhes)
    VALUES (usuario_dono, 'conta aberta', CONCAT('Conta ', NEW.numero_conta, '-', NEW.digito_verificador, ' foi criada.'));
END//

DROP TRIGGER IF EXISTS atualizar_saldo//
CREATE TRIGGER atualizar_saldo AFTER INSERT ON transacao
FOR EACH ROW
BEGIN
    IF NEW.tipo_transacao = 'DEPOSITO' THEN
        UPDATE conta SET saldo = saldo + NEW.valor WHERE id_conta = NEW.id_conta_destino;
    ELSEIF NEW.tipo_transacao IN ('SAQUE', 'TAXA') THEN
        UPDATE conta SET saldo = saldo - NEW.valor WHERE id_conta = NEW.id_conta_origem;
    ELSEIF NEW.tipo_transacao = 'TRANSFERENCIA' THEN
        UPDATE conta SET saldo = saldo - NEW.valor WHERE id_conta = NEW.id_conta_origem;
        UPDATE conta SET saldo = saldo + NEW.valor WHERE id_conta = NEW.id_conta_destino;
    END IF;
END//

DROP TRIGGER IF EXISTS validar_senha//
CREATE TRIGGER validar_senha BEFORE UPDATE ON usuario
FOR EACH ROW
BEGIN
    IF OLD.senha_hash <> NEW.senha_hash THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Senha deve ser atualizada via procedure com validação';
    END IF;
END//

DROP TRIGGER IF EXISTS limite_deposito//
CREATE TRIGGER limite_deposito BEFORE INSERT ON transacao
FOR EACH ROW
BEGIN
    DECLARE total_dia DECIMAL(15,2);
    SELECT COALESCE(SUM(valor), 0) INTO total_dia
    FROM transacao
    WHERE id_conta_origem = NEW.id_conta_origem
      AND tipo_transacao = 'DEPOSITO'
      AND DATE(data_hora) = DATE(NEW.data_hora);
    IF (total_dia + NEW.valor) > 10000 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Limite diário de depósito excedido';
    END IF;
END//

DELIMITER ;
