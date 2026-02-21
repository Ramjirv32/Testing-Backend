# Quick Start: Database Seeding

## 🚀 Quick Commands

### Option 1: Using Bash Script (Recommended)

```bash
cd Backend(Go)
./scripts/seed.sh seed      # Add test data
./scripts/seed.sh clear     # Remove all data
./scripts/seed.sh reset     # Clear and reseed
```

### Option 2: Direct Go Commands

```bash
cd Backend(Go)
go run scripts/seed_data.go    # Seed data
go run scripts/clear_db.go     # Clear database
```

---

## 📊 What Gets Created

### Test Organizer Account
```
Email: organizer@ticpin.com
Password: hashed_password_here
Categories: play, dining, event
```

### Play Venues (2)
- **Elite Badminton Club** - ₹500/slot
  - Location: 12.9716°N, 77.5946°E
  - 8 international standard courts
  
- **Championship Tennis Arena** - ₹800-₹1000/slot
  - Location: 12.9352°N, 77.6365°E
  - Hard and clay courts

### Dining Venues (2)
- **The Golden Plate** (Fine Dining)
  - Rating: 4.8/5
  - Location: 12.9698°N, 77.6109°E
  - Prix fixe & a la carte
  
- **Street Kitchen Dhaba** (Casual)
  - Rating: 4.5/5
  - Location: 12.9716°N, 77.6412°E
  - Street food & dhabha cuisine

### Events (3)
- **Jazz Night** - 15 days from now
  - Artists: Marcus Johnson, Sarah Williams
  - Tickets: ₹500-₹1500
  
- **Comedy Show** - 20 days from now
  - Artist: Aditya Verma
  - Tickets: ₹400-₹800
  
- **Art Exhibition** - 25 days from now (12-day duration)
  - Curator: Priya Sharma
  - Tickets: ₹200-₹500

---

## 📍 Locations

All test data uses real Bangalore coordinates:

```
Whitefield: 12.9698°N, 77.7499°E
Koramangala: 12.9352°N, 77.6365°E
MG Road: 12.9698°N, 77.6109°E
Indiranagar: 12.9716°N, 77.6412°E
Cubbon Park: 12.9716°N, 77.5946°E
```

---

## ✨ Features Included

✅ Real images from Unsplash  
✅ Google Maps integration  
✅ Complete booking settings  
✅ FAQs and T&Cs  
✅ Artist/performer details  
✅ Contact information  
✅ Facility descriptions  
✅ Proper timestamps  

---

## 🔄 Workflow

**First Time Setup:**
```bash
./scripts/seed.sh seed    # Add test data
```

**Start Fresh:**
```bash
./scripts/seed.sh reset   # Delete all and reseed
```

**Clean Database:**
```bash
./scripts/seed.sh clear   # Delete everything
```

---

## 🐛 Troubleshooting

**Problem: "Firestore client is nil"**
- Check `.env` file has correct Firebase credentials
- Verify `config.InitFirebase()` is working

**Problem: Script doesn't run**
- Ensure you're in `Backend(Go)` directory
- Run `chmod +x scripts/seed.sh`

**Problem: No data in Firebase**
- Check Firestore rules allow writes
- Verify database connection in console

---

## 📝 Sample Data Details

### Play Venue Example
```json
{
  "name": "Elite Badminton Club",
  "city": "Bangalore",
  "location": {
    "latitude": 12.9716,
    "longitude": 77.5946,
    "map_url": "https://maps.google.com/?q=12.9716,77.5946"
  },
  "playOptions": [
    {
      "sport": "Badminton",
      "courtType": "Indoor",
      "surface": "Wooden",
      "pricePerSlot": 500
    }
  ],
  "slotSettings": {
    "slotDurationMinutes": 60,
    "openTime": "06:00 AM",
    "closeTime": "10:00 PM"
  }
}
```

### Dining Venue Example
```json
{
  "name": "The Golden Plate",
  "city": "Bangalore",
  "rating": 4.8,
  "location": {
    "latitude": 12.9698,
    "longitude": 77.6109,
    "map_url": "https://maps.google.com/?q=12.9698,77.6109"
  },
  "seatingTypes": [
    {
      "type": "Indoor",
      "totalTables": 20,
      "capacityPerTable": 4
    }
  ],
  "bookingSettings": {
    "advanceBookingDays": 30,
    "averageDiningDurationMinutes": 90
  }
}
```

### Event Example
```json
{
  "title": "Live Jazz Night with The Blue Notes",
  "category": "Music",
  "startDatetime": "2026-03-06T07:00:00Z",
  "endDatetime": "2026-03-06T10:00:00Z",
  "venue": {
    "latitude": 12.9698,
    "longitude": 77.7499,
    "map_url": "https://maps.google.com/?q=12.9698,77.7499"
  },
  "artists": [
    {
      "name": "Marcus Johnson",
      "role": "Lead Saxophone",
      "rating": 4.9,
      "verified": true
    }
  ],
  "tickets": [
    {
      "type": "VIP Seating",
      "price": 1000,
      "available": 35
    }
  ]
}
```

---

## 📚 More Info

See [SEED_DATA_README.md](SEED_DATA_README.md) for complete documentation.
