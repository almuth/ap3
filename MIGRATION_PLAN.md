# AhadPOS 3 to Go REST API Migration Plan

## Executive Summary

This document outlines a comprehensive plan to migrate the AhadPOS 3 system from a PHP/Yii 1.1/MySQL architecture to a modern Go-based REST API with SQLite database. The migration will be executed in phases to minimize disruption and ensure business continuity.

## Current System Overview

### Tech Stack
- **Backend**: PHP 5.6-7.4
- **Framework**: Yii 1.1 (Classic MVC)
- **Database**: MySQL
- **Frontend**: Theme-based HTML/Zurb Foundation 5.5
- **PDF Generation**: mPDF
- **Role-based Access**: CDbAuthManager

### Core Modules
1. **Product/Inventory Management** (Barang, Kategori, Rak, Satuan, Struktur)
2. **Sales Management** (Penjualan, PenjualanDetail, Diskon)
3. **Purchase Management** (Pembelian, PembelianDetail, PO)
4. **Returns Management** (ReturPenjualan, ReturPembelian)
5. **Customer/Member Management** (Profil, Member, Membership)
6. **Supplier Management** (Profil, SupplierBarang)
7. **Cash/Bank Management** (KasBank, Penerimaan, Pengeluaran)
8. **Debts/Credits** (HutangPiutang)
9. **Stock Management** (StockOpname, InventoryBalance)
10. **Reporting** (Multiple report modules)
11. **User/Authentication** (User, AuthAssignment, AuthItem)

### Database Schema
- 136 migrations tracking database evolution
- 103 model classes
- 50+ controllers
- ~50 database tables

## Target Architecture

### Tech Stack
- **Backend**: Go 1.21+
- **Framework**: Chi Router + GORM ORM
- **Database**: SQLite 3
- **Authentication**: JWT (JSON Web Tokens)
- **API Documentation**: Swagger/OpenAPI 3.0
- **Testing**: Go testing framework + testify
- **Deployment**: Docker containers

### Design Principles
1. **RESTful API Design**: Clean, resource-oriented endpoints
2. **Separation of Concerns**: Clear layering (handler → service → repository)
3. **Type Safety**: Leverage Go's strong typing
4. **Performance**: Optimized queries and caching
5. **Scalability**: Stateless API design
6. **Testability**: Comprehensive unit and integration tests
7. **Documentation**: Auto-generated API documentation

## Proposed Go Project Structure

```
ahadpos-go/
├── cmd/
│   ├── api/
│   │   └── main.go                 # API entry point
│   └── migrate/
│       └── main.go                 # Database migration tool
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration management
│   ├── controllers/
│   │   ├── auth.go                 # Authentication handlers
│   │   ├── barang.go               # Product handlers
│   │   ├── penjualan.go            # Sales handlers
│   │   ├── pembelian.go            # Purchase handlers
│   │   ├── kasir.go                # Cashier handlers
│   │   ├── laporan.go              # Report handlers
│   │   └── ...
│   ├── middleware/
│   │   ├── auth.go                 # JWT authentication
│   │   ├── logger.go               # Request logging
│   │   ├── cors.go                 # CORS handling
│   │   └── recovery.go             # Panic recovery
│   ├── models/
│   │   ├── user.go                 # User model
│   │   ├── barang.go               # Product model
│   │   ├── penjualan.go            # Sales model
│   │   ├── pembelian.go            # Purchase model
│   │   └── ...
│   ├── repository/
│   │   ├── user_repository.go      # User data access
│   │   ├── barang_repository.go    # Product data access
│   │   ├── penjualan_repository.go # Sales data access
│   │   └── ...
│   ├── services/
│   │   ├── auth_service.go         # Authentication business logic
│   │   ├── barang_service.go       # Product business logic
│   │   ├── penjualan_service.go    # Sales business logic
│   │   └── ...
│   └── utils/
│       ├── jwt.go                  # JWT utilities
│       ├── validator.go            # Input validation
│       └── response.go             # Response formatting
├── pkg/
│   ├── database/
│   │   ├── database.go             # Database connection
│   │   └── migrations.go           # Migration scripts
│   └── errors/
│       └── errors.go               # Custom error types
├── migrations/
│   ├── 000001_init_schema.sql      # Initial schema
│   ├── 000002_add_users.sql
│   └── ...
├── docs/
│   └── swagger.yaml                # API documentation
├── configs/
│   ├── config.yaml                 # Application config
│   └── config.example.yaml
├── scripts/
│   ├── setup.sh                    # Setup script
│   └── seed_data.sql               # Initial seed data
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## API Endpoint Design

### Authentication Endpoints
```
POST   /api/v1/auth/login          - User login
POST   /api/v1/auth/logout         - User logout
POST   /api/v1/auth/refresh        - Refresh JWT token
GET    /api/v1/auth/me             - Get current user info
```

### Product/Inventory Endpoints
```
GET    /api/v1/barang              - List products (paginated, filtered)
GET    /api/v1/barang/:id          - Get product details
POST   /api/v1/barang              - Create new product
PUT    /api/v1/barang/:id          - Update product
DELETE /api/v1/barang/:id          - Delete product
GET    /api/v1/barang/:id/harga    - Get product prices
GET    /api/v1/barang/:id/stok     - Get product stock
GET    /api/v1/kategori            - List categories
POST   /api/v1/kategori            - Create category
GET    /api/v1/satuan              - List units
POST   /api/v1/satuan              - Create unit
GET    /api/v1/rak                 - List racks
POST   /api/v1/rak                 - Create rack
GET    /api/v1/struktur            - List product structures
```

### Sales Endpoints
```
GET    /api/v1/penjualan           - List sales
GET    /api/v1/penjualan/:id       - Get sale details
POST   /api/v1/penjualan           - Create sale
PUT    /api/v1/penjualan/:id       - Update sale
DELETE /api/v1/penjualan/:id       - Delete sale
GET    /api/v1/penjualan/:id/detail - Get sale items
POST   /api/v1/penjualan/:id/detail - Add item to sale
PUT    /api/v1/penjualan/:id/detail/:detailId - Update item
DELETE /api/v1/penjualan/:id/detail/:detailId - Remove item
POST   /api/v1/penjualan/:id/pay   - Process payment
GET    /api/v1/penjualan/:id/print - Print receipt/invoice
```

### Purchase Endpoints
```
GET    /api/v1/pembelian           - List purchases
GET    /api/v1/pembelian/:id       - Get purchase details
POST   /api/v1/pembelian           - Create purchase
PUT    /api/v1/pembelian/:id       - Update purchase
GET    /api/v1/pembelian/:id/detail - Get purchase items
POST   /api/v1/pembelian/:id/detail - Add item to purchase
GET    /api/v1/po                  - List purchase orders
POST   /api/v1/po                  - Create PO
```

### Customer/Supplier Endpoints
```
GET    /api/v1/profil              - List profiles (customers/suppliers)
GET    /api/v1/profil/:id          - Get profile details
POST   /api/v1/profil              - Create profile
PUT    /api/v1/profil/:id          - Update profile
GET    /api/v1/profil/customer     - List customers only
GET    /api/v1/profil/supplier     - List suppliers only
GET    /api/v1/member              - List members
POST   /api/v1/member              - Register member
GET    /api/v1/member/:id/poin     - Get member points
```

### Cash/Bank Endpoints
```
GET    /api/v1/kasbank             - List cash/bank accounts
POST   /api/v1/kasbank             - Create account
GET    /api/v1/penerimaan          - List receipts
POST   /api/v1/penerimaan          - Create receipt
GET    /api/v1/pengeluaran         - List expenses
POST   /api/v1/pengeluaran         - Create expense
```

### Report Endpoints
```
GET    /api/v1/laporan/harian      - Daily report
GET    /api/v1/laporan/penjualan   - Sales report
GET    /api/v1/laporan/pembelian   - Purchase report
GET    /api/v1/laporan/stok        - Stock report
GET    /api/v1/laporan/laba_rugi   - Profit/loss report
GET    /api/v1/laporan/hutang      - Debt report
GET    /api/v1/laporan/pdf/:type    - Generate PDF report
```

### POS/Cashier Endpoints
```
GET    /api/v1/pos/session         - Get POS session
POST   /api/v1/pos/session         - Start POS session
PUT    /api/v1/pos/session         - Update POS session
POST   /api/v1/pos/scan            - Scan barcode
GET    /api/v1/pos/keranjang       - Get cart items
POST   /api/v1/pos/keranjang       - Add item to cart
PUT    /api/v1/pos/keranjang/:id    - Update cart item
DELETE /api/v1/pos/keranjang/:id    - Remove from cart
POST   /api/v1/pos/checkout        - Process checkout
POST   /api/v1/pos/hold            - Hold transaction
GET    /api/v1/pos/hold            - List held transactions
POST   /api/v1/pos/hold/:id        - Resume held transaction
```

## Database Migration Strategy

### MySQL to SQLite Conversion Challenges
1. **Data Types**: MySQL-specific types need mapping to SQLite
2. **Foreign Keys**: SQLite has different constraint syntax
3. **Indexes**: Index creation syntax differences
4. **Auto-increment**: Different mechanisms for auto-increment
5. **Date/Time Functions**: Different function names and formats
6. **ENUM**: SQLite doesn't have ENUM, use TEXT with CHECK constraint
7. **Full-text Search**: Different FTS implementation

### Type Mapping Guide
```
MySQL                     → SQLite
--------------------------------------------------
TINYINT(1)               → INTEGER (for boolean)
INT(10) UNSIGNED         → INTEGER
BIGINT                   → INTEGER
VARCHAR(30)              → TEXT
DECIMAL(18,2)            → REAL (or store as INTEGER in cents)
TIMESTAMP                → TEXT (ISO 8601 format: "2006-01-02 15:04:05")
DATE                     → TEXT (ISO 8601 format: "2006-01-02")
TEXT                     → TEXT
BLOB                     → BLOB
```

### Migration Steps

#### Phase 1: Schema Analysis
1. Export all MySQL migrations
2. Document all tables with their relationships
3. Create ER diagram of current schema
4. Identify MySQL-specific features that need conversion

#### Phase 2: SQLite Schema Creation
1. Create SQLite-equivalent DDL scripts
2. Handle foreign key constraints properly
3. Create necessary indexes for performance
3. Add triggers for timestamps (since SQLite lacks ON UPDATE CURRENT_TIMESTAMP)

#### Phase 3: Data Migration
1. Export MySQL data to CSV files
2. Transform data types as needed
3. Import into SQLite
4. Validate data integrity
5. Run data consistency checks

#### Phase 4: Validation
1. Compare record counts
2. Verify relationships
3. Test critical queries
4. Performance testing

### Key Schema Considerations

#### Timestamp Handling
SQLite doesn't support `ON UPDATE CURRENT_TIMESTAMP`. Use triggers:
```sql
CREATE TRIGGER update_barang_timestamp
AFTER UPDATE ON barang
BEGIN
    UPDATE barang SET updated_at = datetime('now') WHERE id = NEW.id;
END;
```

#### Decimal Precision
For financial calculations, consider storing amounts as INTEGER (in smallest currency unit) to avoid floating-point precision issues.

#### Foreign Keys
Enable foreign keys in SQLite:
```sql
PRAGMA foreign_keys = ON;
```

## Implementation Phases

### Phase 1: Foundation (Weeks 1-2)
**Objectives**: Set up project infrastructure and core functionality

**Tasks**:
1. Initialize Go project structure
2. Set up database layer with GORM
3. Implement configuration management
4. Create migration framework
5. Implement basic middleware (CORS, logging, recovery)
6. Set up testing framework
7. Create API documentation structure

**Deliverables**:
- Project scaffold
- Database connection working
- Basic middleware stack
- Migration tool functional

### Phase 2: Authentication & User Management (Week 3)
**Objectives**: Implement authentication system

**Tasks**:
1. Implement JWT authentication
2. Create User model and repository
3. Implement login/logout endpoints
4. Add role-based access control
5. Implement password hashing
6. Create user management endpoints

**Deliverables**:
- Working authentication system
- User CRUD operations
- Role-based permissions
- JWT token management

### Phase 3: Product/Inventory Module (Weeks 4-5)
**Objectives**: Migrate core inventory functionality

**Tasks**:
1. Create Product (Barang) model
2. Implement product CRUD operations
3. Create related models (Kategori, Rak, Satuan)
4. Implement barcode scanning
5. Add product search and filtering
6. Implement pricing management
7. Stock management features

**Deliverables**:
- Complete product management API
- Category management
- Unit management
- Rack management
- Stock tracking

### Phase 4: Sales Module (Weeks 6-7)
**Objectives**: Implement sales functionality

**Tasks**:
1. Create Sale (Penjualan) and SaleDetail models
2. Implement sale creation and management
3. Implement discount logic
4. Add payment processing
5. Create receipt/invoice generation
6. Implement returns (ReturPenjualan)
7. Add sales reporting

**Deliverables**:
- Sales API endpoints
- Payment processing
- Receipt generation
- Returns handling
- Sales reports

### Phase 5: Purchase Module (Week 8)
**Objectives**: Implement purchase management

**Tasks**:
1. Create Purchase (Pembelian) and PurchaseDetail models
2. Implement purchase order (PO) system
3. Add supplier management
4. Implement purchase returns
5. Create purchase reporting

**Deliverables**:
- Purchase API endpoints
- Purchase order system
- Supplier management
- Purchase reports

### Phase 6: POS/Cashier Module (Week 9)
**Objectives**: Implement POS functionality

**Tasks**:
1. Create POS session management
2. Implement shopping cart
3. Add barcode scanning integration
4. Implement hold/resume transaction
5. Create checkout process
6. Add cash drawer integration
7. Implement printer integration

**Deliverables**:
- POS API endpoints
- Cart management
- Transaction hold/resume
- Checkout process
- Hardware integration hooks

### Phase 7: Financial Module (Week 10)
**Objectives**: Implement financial tracking

**Tasks**:
1. Create Cash/Bank management
2. Implement income/expense tracking
3. Add debt/credit management (HutangPiutang)
4. Create financial reports
5. Implement accounting integration points

**Deliverables**:
- Cash/Bank API
- Income/expense tracking
- Debt/credit management
- Financial reports

### Phase 8: Member/Customer Module (Week 11)
**Objectives**: Implement customer and member features

**Tasks**:
1. Create customer management
2. Implement membership system
3. Add points/rewards system
4. Create member reporting
5. Implement member-specific discounts

**Deliverables**:
- Customer API
- Membership API
- Points system
- Member reports

### Phase 9: Reporting Module (Week 12)
**Objectives**: Complete reporting functionality

**Tasks**:
1. Implement all report types
2. Add PDF generation
3. Create export functionality (CSV, Excel)
4. Add filtering and date ranges
5. Optimize report queries

**Deliverables**:
- Complete reporting API
- PDF generation
- Data export
- Performance-optimized reports

### Phase 10: Testing & Optimization (Week 13)
**Objectives**: Ensure quality and performance

**Tasks**:
1. Complete unit tests
2. Write integration tests
3. Performance testing
4. Load testing
5. Security audit
6. Code optimization

**Deliverables**:
- Test suite
- Performance benchmarks
- Security report
- Optimized codebase

### Phase 11: Deployment & Documentation (Week 14)
**Objectives**: Prepare for production deployment

**Tasks**:
1. Create Docker configuration
2. Write deployment guide
3. Complete API documentation
4. Create user manual
5. Set up CI/CD pipeline
6. Prepare monitoring setup

**Deliverables**:
- Docker images
- Deployment documentation
- Complete API documentation
- User manual
- CI/CD pipeline

## Technical Considerations

### Go Framework Selection
**Recommended**: Chi Router + GORM
- **Chi**: Lightweight, fast, idiomatic Go router
- **GORM**: Feature-rich ORM with SQLite support
- **Alternatives**: Gin, Echo, Fiber

### Authentication Strategy
**JWT-based authentication**:
- Access tokens: 15 minutes
- Refresh tokens: 7 days
- Stored in HTTP-only cookies
- Role-based claims for authorization

### Validation
**Input validation using validator.v9**:
- Struct tag-based validation
- Custom validators for business rules
- Consistent error responses

### Error Handling
**Custom error types**:
- ValidationError for input errors
- NotFoundError for missing resources
- ConflictError for duplicate/conflicting data
- AuthenticationError for auth failures
- AuthorizationError for permission issues

### Logging
**Structured logging with zerolog**:
- JSON format logs
- Different log levels
- Request ID tracking
- Error stack traces

### Caching Strategy
**Multi-level caching**:
1. In-memory cache for frequently accessed data (products, prices)
2. Redis for session management (optional)
3. Database query caching via GORM

### Performance Optimization
1. Database indexing on frequently queried fields
2. Connection pooling
3. Prepared statements
4. Pagination for large datasets
5. Lazy loading for related data
6. Query optimization

## Data Migration Plan

### Pre-Migration Checklist
- [ ] Backup MySQL database
- [ ] Document all custom queries
- [ ] Identify data inconsistencies
- [ ] Plan for data transformation

### Migration Steps

#### Step 1: Export MySQL Data
```bash
# Export structure
mysqldump --no-data ahadpos3 > schema.sql

# Export data
mysqldump --no-create-info ahadpos3 > data.sql
```

#### Step 2: Transform Schema
- Convert MySQL DDL to SQLite
- Handle data type conversions
- Add necessary triggers
- Create foreign key constraints

#### Step 3: Migrate Data
```bash
# Import into SQLite
sqlite3 ahadpos3.db < transformed_schema.sql
sqlite3 ahadpos3.db < transformed_data.sql
```

#### Step 4: Validate
```sql
-- Check record counts
SELECT 'barang' as table_name, COUNT(*) as count FROM barang
UNION ALL
SELECT 'penjualan', COUNT(*) FROM penjualan
-- ... for all tables
```

### Critical Data Transformations

#### Dates
Convert MySQL timestamps to SQLite datetime format:
```sql
UPDATE barang SET created_at = datetime(created_at);
```

#### Booleans
Convert TINYINT(1) to INTEGER:
```go
type Barang struct {
    Status int `gorm:"type:integer"` // 0 or 1
}
```

#### Decimals
For financial accuracy, store as integers:
```go
type PenjualanDetail struct {
    HargaSatuan int64  `gorm:"type:integer"` // In cents/smallest unit
    Total       int64  `gorm:"type:integer"`
}
```

## Testing Strategy

### Unit Tests
- Test each service in isolation
- Mock repositories
- Test business logic
- Achieve >80% code coverage

### Integration Tests
- Test API endpoints end-to-end
- Use test database
- Test with real scenarios
- Validate database transactions

### API Tests
- Test all endpoints
- Validate response formats
- Test error handling
- Test authentication/authorization

### Performance Tests
- Load testing with realistic data
- Measure response times
- Identify bottlenecks
- Optimize slow queries

## Risk Assessment

### High Risk Items
1. **Data Loss During Migration**: Mitigation with backups and validation
2. **Performance Degradation**: Mitigation with proper indexing and caching
3. **Business Logic Gaps**: Mitigation with comprehensive testing
4. **Security Vulnerabilities**: Mitigation with security audit

### Medium Risk Items
1. **Feature Parity**: Some features may not have direct equivalents
2. **Third-party Integrations**: Hardware drivers, printer integration
3. **Learning Curve**: Team unfamiliarity with Go

### Low Risk Items
1. **UI Changes**: Not affected (API-only migration)
2. **Report Formatting**: Can be adjusted in frontend

## Rollback Plan

If critical issues arise:
1. Maintain PHP version running in parallel
2. Switch DNS/load balancer back to PHP version
3. Investigate and fix issues in Go version
4. Retry migration after fixes

## Success Criteria

1. All core features from PHP version available in Go API
2. API response time < 200ms for 90% of requests
3. 100% data integrity maintained
4. Test coverage > 80%
5. No security vulnerabilities
6. Complete API documentation
7. Successful migration of historical data

## Post-Migration Tasks

1. **Monitoring Setup**
   - Application performance monitoring
   - Error tracking (Sentry)
   - Log aggregation
   - Database performance monitoring

2. **Maintenance**
   - Regular security updates
   - Dependency updates
   - Database optimization
   - Feature enhancements

3. **Training**
   - API documentation for developers
   - User guides for frontend teams
   - Operational runbooks

## Appendix

### A. Go Dependencies
```
github.com/go-chi/chi/v5          # HTTP router
github.com/go-chi/cors             # CORS middleware
gorm.io/gorm                       # ORM
gorm.io/driver/sqlite              # SQLite driver
github.com/golang-jwt/jwt/v5       # JWT library
github.com/go-playground/validator # Validation
github.com/rs/zerolog              # Logging
github.com/spf13/viper             # Configuration
github.com/swaggo/swag             # Swagger documentation
github.com/stretchr/testify        # Testing
```

### B. Migration Scripts Location
- MySQL exports: `scripts/mysql_exports/`
- SQLite schemas: `migrations/`
- Data transformation: `scripts/transform/`
- Validation: `scripts/validate/`

### C. API Response Format
```json
{
  "success": true,
  "data": {
    // Response data
  },
  "error": null,
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}
```

### D. Error Response Format
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "barcode",
        "message": "Barcode is required"
      }
    ]
  }
}
```

## Timeline Summary

| Phase | Duration | Deliverable |
|-------|----------|-------------|
| Foundation | 2 weeks | Project setup, database layer |
| Authentication | 1 week | JWT auth, user management |
| Product/Inventory | 2 weeks | Product management API |
| Sales | 2 weeks | Sales API, payment processing |
| Purchase | 1 week | Purchase API, PO system |
| POS/Cashier | 1 week | POS API, cart, checkout |
| Financial | 1 week | Cash/bank, debt/credit |
| Member/Customer | 1 week | Membership API, points |
| Reporting | 1 week | Complete reporting |
| Testing & Optimization | 1 week | Test suite, performance |
| Deployment & Documentation | 1 week | Production ready |

**Total Timeline**: 14 weeks (3.5 months)

## Conclusion

This migration plan provides a comprehensive roadmap for converting the AhadPOS 3 PHP application to a modern Go-based REST API with SQLite database. The phased approach minimizes risk while ensuring business continuity. Key considerations include data integrity, performance optimization, security, and maintaining feature parity with the existing system.

The Go-based architecture will provide:
- Better performance and scalability
- Type safety and reliability
- Easier deployment and maintenance
- Modern development practices
- Better testing capabilities

Success depends on thorough planning, careful execution, and comprehensive testing at each phase.
