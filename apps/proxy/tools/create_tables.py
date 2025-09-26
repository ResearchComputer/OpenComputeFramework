#!/usr/bin/env python3
"""
Database table creation utility for OCF Proxy Service
Creates necessary tables for API keys, users, and usage tracking
"""

import asyncio
import os
import sys
from pathlib import Path

# Add the proxy directory to the path so we can import our modules
sys.path.insert(0, str(Path(__file__).parent.parent))

# Import database models and configuration from the proxy module
from proxy.database import Base, User, APIKey, APIUsageLog, init_db, DATABASE_URL

# Fix database URL format if needed
database_url = DATABASE_URL or "postgresql://user:password@localhost/ocf_proxy"
if database_url.startswith("postgres://"):
    database_url = database_url.replace("postgres://", "postgresql://", 1)

async def create_tables():
    """Create all database tables"""
    print("Creating database tables...")

    # Use the init_db function from proxy.database
    init_db()

    print("✅ All tables created successfully!")

    # List created tables
    print("\nCreated tables:")
    for table_name in Base.metadata.tables.keys():
        print(f"  - {table_name}")

def main():
    """Main function to run table creation"""
    try:
        # Check if database URL is set
        if not DATABASE_URL or DATABASE_URL == "postgresql://user:password@localhost/ocf_proxy":
            print("⚠️  Warning: Using default database URL.")
            print("   Set PG_URI environment variable to use your actual database.")
            response = input("Continue with default database? (y/N): ").lower()
            if response != 'y':
                print("Exiting...")
                return

        # Create tables
        asyncio.run(create_tables())

        print("\n🎉 Database setup completed successfully!")

    except Exception as e:
        print(f"❌ Error creating tables: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()