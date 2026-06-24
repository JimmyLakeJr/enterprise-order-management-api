# Backend Documentation

Tài liệu chi tiết về project analysis, kiến trúc, và kế hoạch triển khai.

---

## 📚 Documentation Files

### 1. **QUICK_SUMMARY.md** ⭐ START HERE
   - 🎯 Tóm tắt nhanh trạng thái project
   - 🔴 4 vấn đề critical
   - ✅ Những gì đã hoàn thiện
   - 📋 Priority checklist
   - 💡 Design patterns
   - ⏱️ **Đọc: 5 phút**

### 2. **PROJECT_ANALYSIS.md** 📊 COMPREHENSIVE
   - 📋 Tóm tắt tổng quan
   - 🔴 5 lỗi critical
   - ⚠️ 6 cảnh báo
   - 🔨 Danh sách chi tiết chức năng chưa hoàn thiện
   - 📁 Files cần tạo/chỉnh sửa
   - 🎯 Ưu tiên triển khai (4 phase)
   - 🐛 Potential bug sources
   - ⏱️ **Đọc: 15-20 phút**

### 3. **API_ENDPOINTS_PLAN.md** 🔌 IMPLEMENTATION GUIDE
   - 23 endpoints chi tiết (request/response/validation/logic)
   - 📊 Summary table
   - 🎯 5 phases implementation
   - ⏱️ **Đọc: 20-30 phút** (reference during coding)

### 4. **ARCHITECTURE_LAYERS.md** 🏗️ ARCHITECTURE
   - 4-layer architecture overview
   - Layer 1: Handler (8 templates + rules)
   - Layer 2: Service (8 templates + rules)
   - Layer 3: Repository (7 templates + rules)
   - Layer 4: Middleware (3 templates + examples)
   - Complete request flow example
   - Dependency injection pattern
   - Testing strategies
   - ⏱️ **Đọc: 20-25 phút** (reference during coding)

---

## 🚀 QUICK START

### For Project Manager
1. Read: **QUICK_SUMMARY.md** (5 min)
2. Check: **PROJECT_ANALYSIS.md** → "📊 CHECKLIST TRIỂN KHAI" (5 min)
3. Estimate: 27 hours for full implementation

### For Developer
1. Read: **QUICK_SUMMARY.md** (5 min)
2. Read: **ARCHITECTURE_LAYERS.md** (20 min)
3. Reference: **API_ENDPOINTS_PLAN.md** while coding
4. Start with Phase 1: Authentication

### For Code Review
1. Check: **PROJECT_ANALYSIS.md** → "✅ NHỮNG GÌ ĐÃ HOÀN THIỆN"
2. Check: **ARCHITECTURE_LAYERS.md** → "Layer rules"
3. Reference: **API_ENDPOINTS_PLAN.md** for implementation details

---

## 📊 Project Status

| Aspect | Status | Details |
|--------|--------|---------|
| **Database** | ✅ Complete | Schema, indexes, migrations done |
| **Models** | ✅ Complete | All 7 models defined |
| **DTOs** | ✅ Complete | All 5 DTO groups defined |
| **Utilities** | ✅ Complete | JWT, password, validator, response |
| **Config** | ✅ Complete | Environment-based config |
| **Handler** | ⚠️ 8% | Only health endpoint |
| **Service** | ❌ 0% | All empty |
| **Repository** | ❌ 0% | All empty |
| **Middleware** | ❌ 0% | Auth & authorization missing |
| **Routes** | ⚠️ 5% | Only health route |
| **Overall** | ⚠️ 20% | Core foundation ready, layers empty |

---

## 🔴 Critical Issues (MUST FIX FIRST)

1. **Repository Layer KHÔNG tồn tại** - Need all 7 files
2. **Service Layer KHÔNG tồn tại** - Need all 6 files
3. **Handler Layer THIẾU** - Need 5 more files
4. **Middleware KHÔNG tồn tại** - Need JWT auth & authorization
5. **Routes KHÔNG hoàn thiện** - Only health endpoint

---

## 🎯 Implementation Priority

### **PHASE 1: Authentication (Week 1)** 🔐
- JWT Auth Middleware
- Auth Repository (User + RefreshToken)
- Auth Service (register, login, logout, refresh)
- Auth Handler & Routes
- **Test:** Postman register → login → refresh flow

### **PHASE 2: Categories & Products (Week 2)** 📦
- Category Repository, Service, Handler, Routes
- Product Repository, Service, Handler, Routes
- Role Authorization Middleware
- **Test:** Create categories → Create products → List products

### **PHASE 3: Orders (Week 3)** 🛒
- Order Repository (with transaction support)
- Order Service (create with validation, update status)
- Order Handler & Routes
- **Test:** Create order with multiple items → Test stock deduction

### **PHASE 4: Admin & Testing (Week 4)** 🧪
- User Management (list, detail, update, delete)
- API Documentation / Postman collection
- End-to-end testing
- Bug fixes & optimization
- Deployment preparation

---

## 📝 Key Files to Create

### Repository Layer (7 files)
```
user_repository.go (70 lines) - Create, GetByEmail, GetByID, Update, List
product_repository.go (80 lines) - CRUD, List with filters
category_repository.go (50 lines) - CRUD, List
order_repository.go (100 lines) - Create (transaction), List, Update status
refresh_token_repository.go (40 lines) - Create, Get, Revoke
role_repository.go (20 lines) - GetByID, List
```

### Service Layer (6 files)
```
auth_service.go (100 lines) - Register, Login, Logout, Refresh, GetProfile
user_service.go (60 lines) - GetProfile, Update, List, Delete
product_service.go (70 lines) - CRUD, List with filters
category_service.go (40 lines) - CRUD, List
order_service.go (120 lines) - Create (with transaction), List, UpdateStatus
role_service.go (20 lines) - GetByID
```

### Handler Layer (5 files - 1 already exists)
```
auth_handler.go (80 lines) - Register, Login, Logout, Refresh, Profile
user_handler.go (70 lines) - GetProfile, Update, List, Delete
product_handler.go (100 lines) - CRUD, List
category_handler.go (60 lines) - CRUD, List
order_handler.go (100 lines) - Create, List, GetDetail, UpdateStatus
```

### Middleware Layer (3 files)
```
jwt_auth_middleware.go (40 lines) - JWT validation
role_authorization_middleware.go (30 lines) - Role check
cors_middleware.go (20 lines) - CORS setup
```

---

## 🧠 Architecture Overview

```
REQUEST FLOW:
Client → Router → Middleware → Handler → Service → Repository → Database

RESPONSE FLOW:
Database ← Repository ← Service ← Handler ← Router ← Client

LAYER RESPONSIBILITIES:
• Handler: Parse → Validate → Call Service → Return Response
• Service: Validate Rules → Call Repository → Process → Return Result
• Repository: Build SQL → Execute → Map → Return Model
• Middleware: Extract → Validate → Pass to Next
```

---

## ✅ Best Practices Checklist

### Security
- [x] Password hashing (bcrypt)
- [x] JWT token generation
- [ ] JWT validation middleware (TODO)
- [ ] Role authorization middleware (TODO)
- [ ] SQL injection prevention (parameterized queries)
- [ ] CORS configuration
- [ ] No exposed errors to client

### Code Quality
- [x] Clear layer separation
- [x] Consistent error handling
- [x] Validation on input
- [x] DTOs defined
- [ ] Unit tests (TODO)
- [ ] Integration tests (TODO)
- [ ] API documentation (TODO)

### Database
- [x] Schema designed
- [x] Indexes created
- [x] Constraints defined
- [ ] Transaction support in service (TODO)
- [ ] Connection pooling configured
- [ ] Error handling

### Business Logic
- [ ] Order creation with transaction (TODO)
- [ ] Stock deduction validation (TODO)
- [ ] Status transition validation (TODO)
- [ ] Authorization checks (TODO)
- [ ] Email unique check (TODO)
- [ ] Pagination support (TODO)

---

## 📞 Common Questions

### Q1: Where do I start?
**A:** Start with Phase 1 Authentication. Without auth, you can't test other endpoints.

### Q2: Which file is most complex?
**A:** Order Service & Repository (120 lines each) - needs transaction support.

### Q3: How are errors handled?
**A:** Use `util.AppError` with status code → HTTPErrorHandler converts to HTTP response.

### Q4: How is authorization enforced?
**A:** Use `RoleAuthorization` middleware on handler → checks JWT claims.

### Q5: How do I test order creation?
**A:** Use integration test with real database → Check stock reduced after create.

### Q6: Is soft delete or hard delete?
**A:** Soft delete - set `is_active = false`. Do NOT remove from database.

### Q7: How to handle concurrent orders?
**A:** Use database transaction + row-level locking (SELECT ... FOR UPDATE).

---

## 🐛 Common Mistakes to Avoid

1. ❌ Putting business logic in handler
2. ❌ Writing SQL in handler or service
3. ❌ Not using parameterized queries
4. ❌ Not validating permissions
5. ❌ Returning `password_hash` in response
6. ❌ Storing plain refresh tokens in DB
7. ❌ Forgetting transaction in order creation
8. ❌ Not checking if email exists before creating user
9. ❌ Not handling status transition validation
10. ❌ Using `time.Now()` instead of `NOW()` in SQL

---

## 🔗 Related Documentation

**In AGENTS_P1.md:**
- Project requirements
- Stack specifications
- Architecture requirements
- Code rules & conventions
- Security requirements
- Database rules
- Business logic rules
- Permissions matrix

---

## 📞 Support

If you have questions:
1. Check **QUICK_SUMMARY.md** → Common Questions
2. Check **ARCHITECTURE_LAYERS.md** → Layer rules
3. Check **API_ENDPOINTS_PLAN.md** → Endpoint details
4. Check **PROJECT_ANALYSIS.md** → Implementation checklist

---

**Last Updated:** June 20, 2026
**By:** Project Analysis Tool
**Status:** 20% Complete
