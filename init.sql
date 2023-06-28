CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid () UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    middle_name VARCHAR(255) NOT NULL,
    login VARCHAR(255),
    group_number VARCHAR(255) NOT NULL,
    balance float,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
    id UUID DEFAULT gen_random_uuid () UNIQUE,
    type VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    cost  float  NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS histrory (
    user_id UUID NOT NULL,
    task_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

insert into users (id,first_name,last_name,middle_name,login,group_number,balance,created_at,updated_at)
values  ('71ebe4fd-82c9-43b9-ba2e-780bb31df72c','test','test1','test2','Gar7dlcC','test_group',300,current_timestamp,current_timestamp);

insert into tasks (id,type,description,cost,created_at,updated_at)
values ('05b99911-af78-4b50-8ced-800992799a2a','water_level','Дан массив из n неотрицательных целых чисел, представляющий карту высот.
Ширина каждой отметки равна 1.
Вычислите, сколько уровней может быть заполнено водой после дождя.',200,current_timestamp,current_timestamp);