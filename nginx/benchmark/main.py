import argparse
from faker import Faker
from faker.providers import phone_number
from dataclasses import dataclass
from datetime import datetime
from uuid import UUID
import uuid

fake = Faker()
fake.add_provider(phone_number)

@dataclass
class Album:
    id: UUID
    name: str 
    description: str 
    published: bool
    release_date: datetime

@dataclass 
class Comment:
    id: UUID
    user_id: UUID 
    track_id: UUID 
    stars: int 
    text: str

@dataclass 
class Genre:
    id: UUID 
    name: str

@dataclass
class Musician: 
    id: UUID
    name: str 
    email: str 
    password: str 
    salt: str 
    country: str 
    description: str

@dataclass 
class Track:
    id: UUID
    album_id: UUID
    name: str 
    url: str

@dataclass 
class User:
    id: UUID 
    name: str 
    email: str 
    phone: str 
    password: str 
    salt: str 
    country: str

def generate_data(users_cnt, output_file):
    sql_str = ''
    sql_str += "insert into users (id, name, email, phone, password, salt, country) values "
    users = []

    for i in range(users_cnt):
        user = User(id=uuid.uuid4(),
                    name=fake.name() + str(i),
                    email=fake.email() + str(i),
                    phone=fake.phone_number(),
                    password=fake.password(length=40) + str(i),
                    salt=fake.password(length=40) + str(i),
                    country=fake.country())

        while "'" in user.country:
            user.country = fake.country()

        users.append(user)
        sql_str += f"('{user.id}', '{user.name}', '{user.email}', '{user.phone}', "\
                   f"'{user.password}', '{user.salt}', '{user.country}'),\n"

    sql_str = sql_str.strip(',\n') + ';\n\n'
    output_file.write(sql_str)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Generate SQL data for database tables.')
    parser.add_argument('--users', type=int, default=10, help='Number of users to generate')
    parser.add_argument('--tests', type=int, default=50, help='Number of tests to generate')
    parser.add_argument('--output', type=str, default='output.sql', help='Output SQL file')

    args = parser.parse_args()

    with open(args.output, 'w') as file:
        generate_data(args.users, file)

