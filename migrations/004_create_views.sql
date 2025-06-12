-- Views

CREATE OR REPLACE VIEW vw_resumo_contas AS
SELECT c.id_cliente, u.nome, COUNT(co.id_conta) AS total_contas, SUM(co.saldo) AS saldo_total
FROM cliente c
JOIN usuario u ON c.id_usuario = u.id_usuario
JOIN conta co ON c.id_cliente = co.id_cliente
GROUP BY c.id_cliente, u.nome;

CREATE OR REPLACE VIEW vw_movimentacoes_recentes AS
SELECT t.*, c.numero_conta, u.nome AS cliente
FROM transacao t
JOIN conta c ON t.id_conta_origem = c.id_conta
JOIN cliente cl ON c.id_cliente = cl.id_cliente
JOIN usuario u ON cl.id_usuario = u.id_usuario
WHERE t.data_hora >= NOW() - INTERVAL 90 DAY;
