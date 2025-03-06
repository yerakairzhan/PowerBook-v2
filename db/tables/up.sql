Create table bot_settings (
                              id serial primary key ,
                              registration boolean default true
);

Create table users (
                       userid varchar(50) not null primary key ,
                       username varchar(50) not null,
                       registered boolean default false,
                       language VARCHAR(3) DEFAULT 'ru',
                       state VARCHAR(255) NULL,
                       created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'Asia/Almaty')
);

Create table reading_logs (
                              userid varchar(50) not null ,
                              username varchar(50) not null,
                              date DATE NOT NULL DEFAULT CURRENT_DATE,
                              minutes_read INT NOT NULL CHECK (minutes_read >= 0),
                              created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'Asia/Almaty'),
                              primary key (userid, date),
                              foreign key (userid) references users(userid) on delete cascade
)