alter table users add column uuid uuid default gen_random_uuid() not null;
