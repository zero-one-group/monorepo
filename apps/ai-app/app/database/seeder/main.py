import asyncio

from app.core.database import Base, Database
from app.database.seeder.user import seed_users


async def main() -> None:
    await Database.initialize()
    async with Database.engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)

    # Process seed
    async with Database.get_session() as session:
        await seed_users(session)
        await session.commit()

    await Database.dispose()


if __name__ == "__main__":
    asyncio.run(main())
