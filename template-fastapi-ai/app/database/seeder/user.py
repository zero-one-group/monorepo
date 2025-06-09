from app.model.user import User
from sqlalchemy.future import select


async def create_user_if_not_exists(session, data: dict) -> None:
    """
    Insert a User(**data) only if there is no existing user
    with the same username.
    """
    stmt = select(User).filter_by(username=data["username"])
    result = await session.execute(stmt)
    if result.scalar_one_or_none() is None:
        session.add(User(**data))


async def seed_users(session) -> None:
    """
    Seed a fixed list of users, skipping any that already exist.
    """
    users = [
        {
            "full_name": "Adhitya",
            "username": "adhitya",
            "phone_number": "+62812345678",
            "email": "adhitya@email.com",
            "password_hash": "HashedPassword",
        },
        {
            "full_name": "Eric",
            "username": "eric",
            "phone_number": "+6287887887887",
            "email": "eric@email.com",
            "password_hash": "HashedPassword",
        },
    ]

    for u in users:
        await create_user_if_not_exists(session, u)
