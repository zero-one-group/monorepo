# FastAPI AI Project Template

Short brief description about the project.

## Check-In-Dance

### Prerequisites

- Python 3.10+
- [Moonrepo](https://moonrepo.dev/docs/getting-started/installation)
- Docker and Docker Compose (optional, for containerized development)

### Installation

1. Create a new project using this template:

```bash
moon generate template-fastapi-ai
```

2. Set up environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Install dependencies:

```bash
moon ai-app:sync
```

### Available Commands

| Command                                                         | Description                                      |
| --------------------------------------------------------------- | ------------------------------------------------ |
| `moon {{ package_name | kebab_case }}:sync`                     | Library synchronization                          |
| `moon {{ package_name | kebab_case }}:dev`                      | Start development server with hot reload         |
| `moon {{ package_name | kebab_case }}:start`                    | Start production server                          |
| `moon {{ package_name | kebab_case }}:migrate`                  | Run database migrations                          |
| `moon {{ package_name | kebab_case }}:migrate-create -- "name"` | Create a new migration with the specified name   |
| `moon {{ package_name | kebab_case }}:migrate-reset`            | Reset all migrations (downgrade and reapply)     |
| `moon {{ package_name | kebab_case }}:seed`                     | Seed the database with dummy data                |
| `moon {{ package_name | kebab_case }}:check-in-dance`           | Run sync, migrate, and seed in sequence          |

## Development
To get started with this template, we recommend the following:

1. Walk through the “Greeting API” example
    - Begin at the router:

    `app/router/openai.py`

    - Inspect the service layer:

    `app/services/greeting.py`

    - Finally, examine the repository layer:

    `app/repository/openai/greeting.py`

2. Explore core abstractions in the app/core folder

    These utilities and base classes are provided to improve developer experience by reducing boilerplate and enforcing consistent patterns across app.

## Production

### Instrumentation
Tracing is enabled exclusively in the production environment. Set `APP_ENVIRONMENT` to `production` to activate tracing. Alternatively, you may customize the tracing rules in `app/core/trace.py`.

For instructions on customizing span tracing, please refer to the example located at:
- `apps/ai-app/app/repository/openai/greeting.py`

### Error Handling

- We’ve set up a global exception handler in `app/main.py`. It catches every `AppError` raised anywhere in your code and turns it into a structured JSON error response with the correct HTTP status.
- To trigger an error response, simply raise an `AppError` from `app/core/exceptions.py` at any layer (repository, service, or even directly in a route).
  ```python
  from app.core.exceptions import AppError

  # inside repository or service
  if something_went_wrong:
      raise AppError(
          message="Invalid user ID",
          status_code=400,
          code="INVALID_ID",
          data={"user_id": supplied_id}
      )
  ```
- The global handler will produce a response like:
  ```json
  {
    "success": false,
    "message": "Invalid user ID",
    "error_code": "INVALID_ID",
    "data": { "user_id": 123 }
  }
  ```
- No further wiring is needed—just raise `AppError` and FastAPI does the rest.

### Useful Links
- [FastAPI Dependency Injection](https://fastapi.tiangolo.com/tutorial/dependencies/)
    - TL;DR: Declare your dependencies as function parameters, and FastAPI will resolve and inject them for you automatically, no manual wiring required.
    - For more details, see any dependency.py in your layers (e.g. app/repository/dependency.py).
- [FastAPI Deployment](https://fastapi.tiangolo.com/deployment/)
- [FastAPI Handling Errors](https://fastapi.tiangolo.com/tutorial/handling-errors/#install-custom-exception-handlers)
