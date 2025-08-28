package python

const (
	// FastAPI Main App Template
	FastAPIMainPy = `from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.core.config import settings
from app.api.v1.api import api_router

app = FastAPI(
    title="<<TILO:.project_name>> API",
    description="A FastAPI application built with TiLoKit",
    version="1.0.0",
    openapi_url=f"{settings.API_V1_STR}/openapi.json"
)

# Set all CORS enabled origins
if settings.BACKEND_CORS_ORIGINS:
    app.add_middleware(
        CORSMiddleware,
        allow_origins=[str(origin) for origin in settings.BACKEND_CORS_ORIGINS],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

app.include_router(api_router, prefix=settings.API_V1_STR)

@app.get("/")
def read_root():
    return {"message": "Welcome to <<TILO:.project_name>> API"}

@app.get("/health")
def health_check():
    return {"status": "healthy", "service": "<<TILO:.project_name>>"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
`

	// FastAPI Config Template
	FastAPIConfigPy = `import secrets
from typing import Any, Dict, List, Optional, Union
from pydantic import AnyHttpUrl, field_validator
from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    model_config = SettingsConfigDict(case_sensitive=True)

    API_V1_STR: str = "/api/v1"
    SECRET_KEY: str = secrets.token_urlsafe(32)

    # BACKEND_CORS_ORIGINS is a JSON-formatted list of origins
    BACKEND_CORS_ORIGINS: List[AnyHttpUrl] = []

    @field_validator("BACKEND_CORS_ORIGINS", mode="before")
    def assemble_cors_origins(cls, v: Union[str, List[str]]) -> Union[List[str], str]:
        if isinstance(v, str) and not v.startswith("["):
            return [i.strip() for i in v.split(",")]
        elif isinstance(v, (list, str)):
            return v
        raise ValueError(v)

    PROJECT_NAME: str = "<<TILO:.project_name>>"

    SQLALCHEMY_DATABASE_URI: Optional[str] = "sqlite:///./sql_app.db"

settings = Settings()
`

	// FastAPI Dependencies Template
	FastAPIDepsPy = `from typing import Generator
from sqlalchemy.orm import Session
from app.db.session import SessionLocal

def get_db() -> Generator:
    try:
        db = SessionLocal()
        yield db
    finally:
        db.close()
`

	// FastAPI Database Session Template
	FastAPISessionPy = `from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from app.core.config import settings

engine = create_engine(settings.SQLALCHEMY_DATABASE_URI, pool_pre_ping=True)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base = declarative_base()
`

	// FastAPI Item Model Template
	FastAPIItemModelPy = `from sqlalchemy import Boolean, Column, Integer, String
from app.db.session import Base

class Item(Base):
    __tablename__ = "items"

    id = Column(Integer, primary_key=True, index=True)
    title = Column(String, index=True)
    description = Column(String, index=True)
    is_active = Column(Boolean, default=True)
`

	// FastAPI User Model Template
	FastAPIUserModelPy = `from sqlalchemy import Boolean, Column, Integer, String
from app.db.session import Base

class User(Base):
    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True)
    email = Column(String, unique=True, index=True)
    hashed_password = Column(String)
    is_active = Column(Boolean, default=True)
`

	// FastAPI Item Schema Template
	FastAPIItemSchemaPy = `from typing import Optional
from pydantic import BaseModel, ConfigDict

class ItemBase(BaseModel):
    title: Optional[str] = None
    description: Optional[str] = None

class ItemCreate(ItemBase):
    title: str

class ItemUpdate(ItemBase):
    pass

class ItemInDBBase(ItemBase):
    model_config = ConfigDict(from_attributes=True)
    id: Optional[int] = None

class Item(ItemInDBBase):
    pass
`

	// FastAPI User Schema Template
	FastAPIUserSchemaPy = `from typing import Optional
from pydantic import BaseModel, EmailStr, ConfigDict

class UserBase(BaseModel):
    email: Optional[EmailStr] = None
    is_active: Optional[bool] = True

class UserCreate(UserBase):
    email: EmailStr
    password: str

class UserUpdate(UserBase):
    password: Optional[str] = None

class UserInDBBase(UserBase):
    model_config = ConfigDict(from_attributes=True)
    id: Optional[int] = None

class User(UserInDBBase):
    pass
`

	// FastAPI API Router Template
	FastAPIApiPy = `from fastapi import APIRouter
from app.api.v1.endpoints import items, users

api_router = APIRouter()
api_router.include_router(items.router, prefix="/items", tags=["items"])
api_router.include_router(users.router, prefix="/users", tags=["users"])
`

	// FastAPI Items Endpoint Template
	FastAPIItemsPy = `from typing import Any, List
from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from app.api import deps
from app.models.item import Item
from app.schemas.item import Item as ItemSchema, ItemCreate

router = APIRouter()

@router.get("/", response_model=List[ItemSchema])
def read_items(
    db: Session = Depends(deps.get_db),
    skip: int = 0,
    limit: int = 100,
) -> Any:
    items = db.query(Item).offset(skip).limit(limit).all()
    return items

@router.post("/", response_model=ItemSchema)
def create_item(
    *,
    db: Session = Depends(deps.get_db),
    item_in: ItemCreate,
) -> Any:
    item = Item(**item_in.model_dump())
    db.add(item)
    db.commit()
    db.refresh(item)
    return item
`

	// FastAPI Users Endpoint Template
	FastAPIUsersPy = `from typing import Any, List
from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from app.api import deps
from app.models.user import User
from app.schemas.user import User as UserSchema

router = APIRouter()

@router.get("/", response_model=List[UserSchema])
def read_users(
    db: Session = Depends(deps.get_db),
    skip: int = 0,
    limit: int = 100,
) -> Any:
    users = db.query(User).offset(skip).limit(limit).all()
    return users
`

	// FastAPI PyProject Template
	FastAPIPyprojectToml = `[tool.poetry]
name = "<<TILO:.project_name>>"
version = "0.1.0"
description = "FastAPI web application"

[tool.poetry.dependencies]
python = "^<<TILO:.python_version>>"
fastapi = "^<<TILO:.fastapi_version>>"
uvicorn = {extras = ["standard"], version = "^0.24.0"}
pydantic = {extras = ["email"], version = "^2.5.0"}
pydantic-settings = "^2.2.1"
sqlalchemy = "^2.0.23"
alembic = "^1.13.0"
python-jose = {extras = ["cryptography"], version = "^3.3.0"}
passlib = {extras = ["bcrypt"], version = "^1.7.4"}
python-multipart = "^0.0.6"
psycopg2-binary = "^2.9.9"

[tool.poetry.group.dev.dependencies]
pytest = "^7.4.0"
pytest-asyncio = "^0.21.0"
httpx = "^0.25.0"
black = "^23.7.0"
isort = "^5.12.0"
mypy = "^1.7.0"

[build-system]
requires = ["poetry-core"]
build-backend = "poetry.core.masonry.api"
`

	// FastAPI Requirements Template
	FastAPIRequirementsTxt = `fastapi==<<TILO:.fastapi_version>>
uvicorn[standard]==0.24.0
pydantic[email]==2.5.0
pydantic-settings==2.2.1
sqlalchemy==2.0.23
alembic==1.13.0
python-jose[cryptography]==3.3.0
passlib[bcrypt]==1.7.4
python-multipart==0.0.6
psycopg2-binary==2.9.9
`

	// FastAPI Pipenv Template
	FastAPIPipfile = `[[source]]
url = "https://pypi.org/simple"
verify_ssl = true
name = "pypi"

[packages]
fastapi = "<<TILO:.fastapi_version>>"
uvicorn = {extras = ["standard"], version = "0.24.0"}
pydantic = {extras = ["email"], version = "2.5.0"}
pydantic-settings = "==2.2.1"
sqlalchemy = "==2.0.23"
alembic = "==1.13.0"
python-jose = {extras = ["cryptography"], version = "3.3.0"}
passlib = {extras = ["bcrypt"], version = "1.7.4"}
python-multipart = "==0.0.6"
psycopg2-binary = "==2.9.9"

[dev-packages]
pytest = "7.4.0"
httpx = "0.25.2"

[requires]
python_version = "<<TILO:.python_version>>"
`

	// FastAPI Conda Environment Template
	FastAPIEnvironmentYml = `name: <<TILO:.project_name>>
channels:
  - conda-forge
  - defaults
dependencies:
  - python=<<TILO:.python_version>>
  - pip
  - pip:
    - fastapi==<<TILO:.fastapi_version>>
    - uvicorn[standard]==0.24.0
    - pydantic[email]==2.5.0
    - pydantic-settings==2.2.1
    - sqlalchemy==2.0.23
    - alembic==1.13.0
    - python-jose[cryptography]==3.3.0
    - passlib[bcrypt]==1.7.4
    - python-multipart==0.0.6
    - psycopg2-binary==2.9.9
    - pytest==7.4.0
    - httpx==0.25.2
`

	// FastAPI README Template
	FastAPIReadmeMd = `# <<TILO:.project_name>>

A FastAPI application built with TiLoKit.

## Setup

1. Create virtual environment:
` + "```bash" + `
python -m venv venv
source venv/bin/activate
` + "```" + `

## 2. Install dependencies (choose one):
` + "```bash" + `
# pip
pip install -r requirements.txt
# poetry
poetry install
# pipenv
pipenv install
# conda
conda env create -f environment.yml && conda activate <<TILO:.project_name>>
` + "```" + `

## 3. Run the application:
` + "```bash" + `
uvicorn main:app --reload
` + "```" + `

Visit http://localhost:8000 to see your application.
Visit http://localhost:8000/docs for interactive API documentation.

## Features

- FastAPI framework
- Automatic API documentation
- Pydantic models
- SQLAlchemy ORM
- Authentication ready
- CORS enabled
`

	// FastAPI Test Configuration Template
	FastAPIConftestPy = `import pytest
from fastapi.testclient import TestClient
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from app.db.session import Base
from app.api.deps import get_db
from main import app

SQLALCHEMY_DATABASE_URL = "sqlite:///./test.db"

engine = create_engine(
    SQLALCHEMY_DATABASE_URL, connect_args={"check_same_thread": False}
)
TestingSessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base.metadata.create_all(bind=engine)

def override_get_db():
    try:
        db = TestingSessionLocal()
        yield db
    finally:
        db.close()

app.dependency_overrides[get_db] = override_get_db

@pytest.fixture
def client():
    return TestClient(app)
`

	// FastAPI Test Main Template
	FastAPITestMainPy = `def test_read_main(client):
    response = client.get("/")
    assert response.status_code == 200
    assert "message" in response.json()

def test_health_check(client):
    response = client.get("/health")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "healthy"

def test_read_items(client):
    response = client.get("/api/v1/items/")
    assert response.status_code == 200
    assert isinstance(response.json(), list)
`
)
