DO $$
	BEGIN
    	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'application_decision_enum') THEN
        	CREATE TYPE application_decision_enum AS ENUM ('Pending', 'Accept', 'Reject');
    	END IF;
	END
$$;

create table if not exists "applications" (
  "id" serial not null,
  "userid" uuid not null,
  "age" smallint not null check (age BETWEEN 1 AND 80),
  "about" varchar(4096) not null,
  "join_reason" varchar(1024) not null,
  "inviter" varchar(256) null,
  "submitted_at" TIMESTAMPTZ not null default NOW(),
  "ai_decision" application_decision_enum not null default 'Pending',
  "admin_decision" application_decision_enum not null default 'Pending',
  "ai_answer" VARCHAR(1024),
  "ai_comment" VARCHAR(4096),
  constraint "users_pkey" primary key ("id")
);

comment on column "applications"."userid" is 'UUID юзера';
comment on column "applications"."age" is 'Возраст игрока 1-80 лет';
comment on column "applications"."about" is 'Поле "О себе" длинной не более 4096 символов';
comment on column "applications"."join_reason" is 'Поле "Почему хотите вступить" длинной не более 1024 символов';
comment on column "applications"."inviter" is 'Поле "Кто вас приласил" длинной не более 256 символов, опционально и может быть null';
comment on column "applications"."submitted_at" is 'Дата и время когда была создана заявка';
comment on column "applications"."ai_decision" is 'Мнение ИИ, принят или нет';
comment on column "applications"."admin_decision" is 'Мнение админов принят или нет';
comment on column "applications"."ai_answer" is 'Ответ ИИ самому игроку';
comment on column "applications"."ai_comment" is 'Пояснение ИИ к своему ответу';
