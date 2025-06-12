-- Functions and stored procedures

DELIMITER //

DROP FUNCTION IF EXISTS calcular_dv//
CREATE FUNCTION calcular_dv(num VARCHAR(20)) RETURNS CHAR(1)
DETERMINISTIC
BEGIN
    DECLARE sum INT DEFAULT 0;
    DECLARE len INT DEFAULT CHAR_LENGTH(num);
    DECLARE i INT DEFAULT 0;
    DECLARE digit INT;
    DECLARE double_digit INT;
    DECLARE parity INT;
    SET parity = len % 2;
    
    WHILE i < len DO
        SET digit = CAST(SUBSTRING(num, i+1, 1) AS UNSIGNED);
        IF (i % 2) = parity THEN
            SET double_digit = digit * 2;
            IF double_digit > 9 THEN
                SET double_digit = double_digit - 9;
            END IF;
            SET sum = sum + double_digit;
        ELSE
            SET sum = sum + digit;
        END IF;
        SET i = i + 1;
    END WHILE;
    RETURN CAST((10 - (sum % 10)) % 10 AS CHAR);
END//

DROP PROCEDURE IF EXISTS gerar_otp//
CREATE PROCEDURE gerar_otp(IN p_id_usuario INT)
BEGIN
    DECLARE novo_otp VARCHAR(6);
    SET novo_otp = LPAD(FLOOR(RAND() * 1000000), 6, '0');
    UPDATE usuario SET otp_ativo = novo_otp, otp_expiracao = NOW() + INTERVAL 5 MINUTE
    WHERE id_usuario = p_id_usuario;
    SELECT novo_otp;
END//

DROP PROCEDURE IF EXISTS calcular_score_credito//
CREATE PROCEDURE calcular_score_credito(IN p_id_cliente INT)
BEGIN
    DECLARE total_trans DECIMAL(15,2);
    DECLARE media_trans DECIMAL(15,2);
    SELECT COALESCE(SUM(valor), 0), COALESCE(AVG(valor), 0) INTO total_trans, media_trans
    FROM transacao t
    JOIN conta c ON t.id_conta_origem = c.id_conta
    WHERE c.id_cliente = p_id_cliente AND t.tipo_transacao IN ('DEPOSITO', 'SAQUE');
    UPDATE cliente SET score_credito = LEAST(100, (total_trans / 1000) + (media_trans / 100))
    WHERE id_cliente = p_id_cliente;
END//

DELIMITER ;
