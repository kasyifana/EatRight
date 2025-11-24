# EatRight API Endpoints Documentation

Base URL: `http://localhost:8080/api`

## Authentication

All protected endpoints require JWT token in header:
```
Authorization: Bearer <jwt-token>
```

---

## Endpoints List

### 🔐 Authentication

#### Verify Supabase Token
```
POST /api/auth/verify
```

**Request Body:**
```json
{
  "supabase_token": "string"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Authentication successful",
  "data": {
    "token": "jwt-token-here",
    "user": {
      "id": "uuid",
      "name": "string",
      "email": "string",
      "role": "user" | "restaurant",
      "created_at": "timestamp"
    }
  }
}
```

---

### 👤 Users

#### Get Current User Profile
```
GET /api/users/me
```
**Auth:** Required  
**Response:**
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": "uuid",
    "name": "string",
    "email": "string",
    "role": "user" | "restaurant",
    "created_at": "timestamp"
  }
}
```

---

### 🏪 Restaurants

#### List Restaurants
```
GET /api/restaurants?lat=<number>&lng=<number>&distance=<number>
```
**Auth:** Public  
**Query Params:**
- `lat` (optional): latitude for nearby search
- `lng` (optional): longitude for nearby search  
- `distance` (optional): max distance in km (default: 10)

**Response:**
```json
{
  "success": true,
  "message": "Restaurants retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "owner_id": "uuid",
      "name": "string",
      "address": "string",
      "lat": 0.0,
      "lng": 0.0,
      "closing_time": "HH:MM:SS",
      "created_at": "timestamp"
    }
  ]
}
```

#### Get Restaurant Detail
```
GET /api/restaurants/:id
```
**Auth:** Public  
**Response:**
```json
{
  "success": true,
  "message": "Restaurant retrieved successfully",
  "data": {
    "id": "uuid",
    "owner_id": "uuid",
    "name": "string",
    "address": "string",
    "lat": 0.0,
    "lng": 0.0,
    "closing_time": "HH:MM:SS",
    "created_at": "timestamp",
    "listings": [...]
  }
}
```

#### Create Restaurant
```
POST /api/restaurants
```
**Auth:** Required (Restaurant role only)  
**Request Body:**
```json
{
  "name": "string",
  "address": "string",
  "lat": 0.0,
  "lng": 0.0,
  "closing_time": "HH:MM:SS"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Restaurant created successfully",
  "data": { ... }
}
```

---

### 📦 Listings

#### List All Active Listings
```
GET /api/listings
```
**Auth:** Public  
**Response:**
```json
{
  "success": true,
  "message": "Listings retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "restaurant_id": "uuid",
      "type": "mystery_box" | "reveal",
      "name": "string | null",
      "description": "string",
      "price": 0,
      "stock": 0,
      "photo_url": "string",
      "pickup_time": "HH:MM:SS",
      "is_active": true,
      "created_at": "timestamp"
    }
  ]
}
```

#### Get Listing Detail
```
GET /api/listings/:id
```
**Auth:** Public  
**Response:**
```json
{
  "success": true,
  "message": "Listing retrieved successfully",
  "data": { ... }
}
```

#### Create Listing
```
POST /api/restaurants/:id/listings
```
**Auth:** Required (Restaurant role only)  
**Request Body:**
```json
{
  "type": "mystery_box" | "reveal",
  "name": "string | null",
  "description": "string",
  "price": 0,
  "stock": 0,
  "photo_url": "string",
  "pickup_time": "HH:MM:SS"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Listing created successfully",
  "data": { ... }
}
```

#### Update Listing Stock
```
PATCH /api/listings/:id/stock
```
**Auth:** Required (Restaurant role only)  
**Request Body:**
```json
{
  "quantity": 0  // positive to add, negative to reduce
}
```

**Response:**
```json
{
  "success": true,
  "message": "Stock updated successfully",
  "data": null
}
```

#### Toggle Listing Status
```
PATCH /api/listings/:id/status
```
**Auth:** Required (Restaurant role only)  
**Request Body:**
```json
{
  "is_active": true | false
}
```

**Response:**
```json
{
  "success": true,
  "message": "Status updated successfully",
  "data": null
}
```

---

### 🛒 Orders

#### Create Order
```
POST /api/orders
```
**Auth:** Required  
**Request Body:**
```json
{
  "listing_id": "uuid",
  "qty": 0
}
```

**Response:**
```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "listing_id": "uuid",
    "qty": 0,
    "total_price": 0,
    "status": "pending",
    "created_at": "timestamp"
  }
}
```

#### Get My Orders
```
GET /api/orders/me
```
**Auth:** Required  
**Response:**
```json
{
  "success": true,
  "message": "Orders retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "listing_id": "uuid",
      "qty": 0,
      "total_price": 0,
      "status": "pending" | "ready" | "completed" | "cancelled",
      "created_at": "timestamp",
      "listing": { ... }
    }
  ]
}
```

#### Update Order Status
```
PATCH /api/orders/:id/status
```
**Auth:** Required (Restaurant role only)  
**Request Body:**
```json
{
  "status": "pending" | "ready" | "completed" | "cancelled"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Order status updated successfully",
  "data": null
}
```

---

## Error Response Format

All errors follow this format:
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

Common HTTP Status Codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation error)
- `401` - Unauthorized (missing/invalid token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `500` - Internal Server Error

---

## Angular Integration Example

```typescript
// api.service.ts
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private baseUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  // Get auth headers
  private getHeaders(): HttpHeaders {
    const token = localStorage.getItem('jwt_token');
    return new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': token ? `Bearer ${token}` : ''
    });
  }

  // Restaurants
  getRestaurants(lat?: number, lng?: number, distance?: number): Observable<any> {
    let url = `${this.baseUrl}/restaurants`;
    const params = new URLSearchParams();
    if (lat) params.append('lat', lat.toString());
    if (lng) params.append('lng', lng.toString());
    if (distance) params.append('distance', distance.toString());
    
    return this.http.get(`${url}?${params.toString()}`);
  }

  getRestaurant(id: string): Observable<any> {
    return this.http.get(`${this.baseUrl}/restaurants/${id}`);
  }

  // Listings
  getListings(): Observable<any> {
    return this.http.get(`${this.baseUrl}/listings`);
  }

  getListing(id: string): Observable<any> {
    return this.http.get(`${this.baseUrl}/listings/${id}`);
  }

  // Orders (protected)
  createOrder(listingId: string, qty: number): Observable<any> {
    return this.http.post(
      `${this.baseUrl}/orders`,
      { listing_id: listingId, qty },
      { headers: this.getHeaders() }
    );
  }

  getMyOrders(): Observable<any> {
    return this.http.get(`${this.baseUrl}/orders/me`, {
      headers: this.getHeaders()
    });
  }

  // Auth
  verifyToken(supabaseToken: string): Observable<any> {
    return this.http.post(`${this.baseUrl}/auth/verify`, {
      supabase_token: supabaseToken
    });
  }
}
```

---

## Testing with cURL

```bash
# Public endpoints (no auth needed)
curl http://localhost:8080/api/restaurants
curl http://localhost:8080/api/listings

# Protected endpoints (need JWT)
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/users/me

curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"listing_id":"uuid","qty":1}' \
  http://localhost:8080/api/orders
```
