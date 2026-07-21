<div align="center">

# ERP API

**An Enterprise Resource Planning (ERP) backend built with Go**

Inspired by enterprise systems such as **SAP ERP**. This project focuses on modeling real-world business processes using modern backend technologies.

![Go](https://img.shields.io/badge/Go-1.26.4-blue)
![Echo](https://img.shields.io/badge/Echo-v5-green)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-blue)
![Redis](https://img.shields.io/badge/Redis-Cache-red)
![GORM](https://img.shields.io/badge/GORM-ORM-orange)

</div>

---

# Overview

ERP API is an educational enterprise backend project that aims to understand how modern ERP systems manage business operations from **master data creation** to **procurement**, **inventory**, **warehouse**, **sales**, and **financial processes**.

Unlike a simple CRUD application, this project focuses on modeling **enterprise business workflows** and implementing concepts commonly found in ERP platforms such as SAP ERP.

The project is currently developed as a **Modular Monolith** with a long-term roadmap toward **Microservices Architecture**.

---

# Project Goals

- Learn Enterprise Software Development
- Understand ERP Architecture
- Study SAP-inspired Business Processes
- Practice Clean Architecture
- Build scalable backend applications using Go
- Model real-world Supply Chain Management workflows

---

# Business Process

```
                  Material Master
                         │
                         ▼
                  Vendor Management
                         │
                         ▼
             Purchase Requisition
                         │
                         ▼
          Request for Quotation (RFQ)
                         │
                         ▼
                 Purchase Order
                         │
                         ▼
                 Goods Receipt
                         │
                         ▼
              Warehouse Inventory
                         │
                         ▼
                 Sales Order
                         │
                         ▼
                    Delivery
                         │
                         ▼
                     Invoice
```

---

# Features

## Master Data

- 

---

## Procurement

- Purchase Requisition *(Planned)*
- Request for Quotation *(Planned)*
- Purchase Order *(Planned)*
- Goods Receipt *(Planned)*

---

## Inventory

- Stock Management *(Planned)*
- Goods Issue *(Planned)*
- Stock Transfer *(Planned)*
- Stock Adjustment *(Planned)*

---

## Warehouse

- Warehouse
- Storage Location
- Bin Management *(Planned)*

---

## Sales

- Customer *(Planned)*
- Sales Organization *(Planned)*
- Distribution Channel *(Planned)*
- Sales Order *(Planned)*
- Delivery *(Planned)*
- Billing *(Planned)*

---

## Finance

- Accounting *(Planned)*
- General Ledger *(Planned)*
- Invoice *(Planned)*

---

# Technology Stack

| Layer | Technology |
|--------|------------|
| Language | Go |
| Framework | Echo |
| ORM | GORM |
| Database | PostgreSQL |
| Cache | Redis |
| API | REST API |
| Authentication | Session Authentication |
| Documentation | Swagger |
| Container | Docker |

---

# Architecture

Current implementation follows a **Modular Monolithic Architecture**.

```
Client
    │
REST API
    │
Controller
    │
Service
    │
Repository
    │
PostgreSQL
```

Each module is isolated to simplify future migration into microservices.

---

# Design Principles

- Clean Architecture
- Repository Pattern
- Dependency Injection
- RESTful API
- Modular Design
- Domain-Oriented Design
- SOLID Principles

---

# Learning Objectives

This project is used to deepen understanding of:

- Enterprise Resource Planning (ERP)
- Material Management
- Supply Chain Management
- Procurement Process
- Inventory Management
- Warehouse Management
- Business Process Modeling
- Database Normalization
- Enterprise Software Design

---

# Roadmap

## Phase 1 — Foundation

- [x] Material Master
- [x] Material Type
- [x] Material Group
- [x] Material Unit
- [x] Company
- [x] Warehouse

---

## Phase 2 — Procurement

- [ ] Vendor Master
- [ ] Purchase Requisition
- [ ] RFQ
- [ ] Purchase Order
- [ ] Goods Receipt

---

## Phase 3 — Inventory

- [ ] Inventory
- [ ] Stock Movement
- [ ] Warehouse Transfer
- [ ] Cycle Counting

---

## Phase 4 — Sales

- [ ] Customer
- [ ] Sales Order
- [ ] Delivery
- [ ] Billing

---

## Phase 5 — Finance

- [ ] Accounting
- [ ] General Ledger
- [ ] Financial Report

---

# Future Vision

The long-term goal of this project is to evolve into a complete enterprise ERP platform capable of supporting medium-sized businesses while serving as a learning platform for enterprise software engineering.

---

# Inspiration

This project is inspired by concepts from:

- SAP ERP
- SAP S/4HANA

This project is developed independently for educational purposes and does not contain proprietary implementations from any commercial ERP vendor.
