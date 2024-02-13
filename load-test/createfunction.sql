CREATE OR REPLACE FUNCTION createtransaction(
    IN idUser integer,
    IN value integer,
    IN description varchar(10)
) RETURNS RECORD AS $$
DECLARE
userfound user_models%rowtype;
    ret RECORD;
BEGIN
SELECT * FROM user_models
    INTO userfound
WHERE id = idUser;

IF not found THEN
        --raise notice'Id user % not found.', idUser;
select -1 into ret;
RETURN ret;
END IF;

    --raise notice'service by user %.', idUser;
INSERT INTO models (value, description, created_at, user_id)
VALUES (value, description, now() at time zone 'utc', idUser);
UPDATE user_models
SET balance = balance + value
WHERE id = idUser AND (value > 0 OR balance + value >= "limit")
    RETURNING balance, "limit"
INTO ret;
raise notice'Ret: %', ret;
    IF ret."limit" is NULL THEN
        --raise notice'Id  user % not found.', idUser;
select -2 into ret;
END IF;
RETURN ret;
END;$$ LANGUAGE plpgsql;


select createtransaction(1, -1000, 'foo');


SELECT user_models.balance, user_models.limit, models.value, models.description, models.created_at
FROM models
          RIGHT JOIN user_models on user_models.id = models.user_id
WHERE user_models.id = 1
order by models.created_at desc
limit 10;

select * from user_models;