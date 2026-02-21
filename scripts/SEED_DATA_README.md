# Database Seeding Script

This directory contains scripts to manage your TicPin database.

## Available Scripts

### 1. `seed_data.go` - Add Real Test Data

Seeds your Firebase database with realistic test data for:
- **Play Venues**: Badminton courts, Tennis arenas with real locations in Bangalore
- **Dining Venues**: Restaurants with real addresses, menus, and booking settings
- **Events**: Live music, comedy shows, art exhibitions with artist details
- **Organizer Account**: Creates an admin user to manage all seeded content

#### Data Included:
- Real Bangalore locations with actual GPS coordinates
- Google Maps links for each venue
- Professional images from Unsplash
- Complete booking settings and configurations
- FAQs and terms & conditions
- Contact details and facilities information

#### How to Run:

1. **Navigate to Backend directory:**
   ```bash
   cd Backend(Go)
   ```

2. **Run the seed script:**
   ```bash
   go run scripts/seed_data.go
   ```

#### What Gets Created:

**Organizer Account:**
- Email: `organizer@ticpin.com`
- Categories: Play, Dining, Events

**Play Venues (2 items):**
1. Elite Badminton Club (Bangalore, 12.9716, 77.5946)
   - Badminton courts with professional coaching
   - 60-minute slots at ₹500

2. Championship Tennis Arena (Bangalore, 12.9352, 77.6365)
   - Hard and clay courts
   - 60-minute slots at ₹800-₹1000

**Dining Venues (2 items):**
1. The Golden Plate (Bangalore, 12.9698, 77.6109)
   - Fine dining Indian restaurant
   - Rating: 4.8/5
   - 90-minute dining slots

2. Street Kitchen Dhaba (Bangalore, 12.9716, 77.6412)
   - Authentic street food
   - Rating: 4.5/5
   - 45-minute casual dining

**Events (3 items):**
1. Live Jazz Night with The Blue Notes
   - **When:** 15 days from now
   - **Where:** The Jazz Lounge (12.9698, 77.7499)
   - **Price:** ₹500-₹1500
   - **Artists:** Marcus Johnson (Saxophone), Sarah Williams (Piano)

2. Comedy Night: Stand-up Comedy Extravaganza
   - **When:** 20 days from now
   - **Where:** Comedy Central Theater (12.9387, 77.6821)
   - **Price:** ₹400-₹800
   - **Artist:** Aditya Verma

3. Art Exhibition: Contemporary Indian Artists
   - **When:** 25 days from now (ongoing for 12 days)
   - **Where:** Bangalore Art Museum (12.9716, 77.5946)
   - **Price:** ₹200-₹500
   - **Curator:** Priya Sharma

---

### 2. `clear_db.go` - Clear All Data

Removes all collections from the database. Use with caution!

```bash
go run scripts/clear_db.go
```

---

## Workflow: Seeding New Data

**Step 1:** Clear existing data (optional)
```bash
go run scripts/clear_db.go
```

**Step 2:** Seed new data
```bash
go run scripts/seed_data.go
```

**Step 3:** Verify in Firebase Console
- Go to your Firebase project Firestore
- Check the collections: `play_venues`, `dining_venues`, `events`, `users`

---

## Locations in Test Data

All venues are located in Bangalore with real GPS coordinates:

| Venue | Latitude | Longitude | Map URL |
|-------|----------|-----------|---------|
| Elite Badminton Club | 12.9716 | 77.5946 | [Map](https://maps.google.com/?q=12.9716,77.5946) |
| Championship Tennis Arena | 12.9352 | 77.6365 | [Map](https://maps.google.com/?q=12.9352,77.6365) |
| The Golden Plate | 12.9698 | 77.6109 | [Map](https://maps.google.com/?q=12.9698,77.6109) |
| Street Kitchen Dhaba | 12.9716 | 77.6412 | [Map](https://maps.google.com/?q=12.9716,77.6412) |
| The Jazz Lounge | 12.9698 | 77.7499 | [Map](https://maps.google.com/?q=12.9698,77.7499) |
| Comedy Central Theater | 12.9387 | 77.6821 | [Map](https://maps.google.com/?q=12.9387,77.6821) |
| Bangalore Art Museum | 12.9716 | 77.5946 | [Map](https://maps.google.com/?q=12.9716,77.5946) |

---

## Customizing the Data

To add more venues or modify existing ones:

1. Edit `seed_data.go`
2. Add new venue entries to the respective slice
3. Save and run: `go run scripts/seed_data.go`

---

## Notes

- Images are from Unsplash (free to use)
- All contact emails and phone numbers are test data
- GPS coordinates are real Bangalore locations
- Events are scheduled 15, 20, and 25 days from current date
- Tickets have realistic pricing and availability
- All records include proper timestamps

---

## Troubleshooting

**Error: "Firestore client is nil"**
- Ensure Firebase credentials are properly set in `.env`
- Check that `config.InitFirebase()` completes successfully

**No data appears in Firestore**
- Verify Firebase project credentials
- Check Firestore database rules allow writes
- Review logs for specific error messages

**Duplicate data after re-running**
- Run `clear_db.go` first to remove old data
- Then run `seed_data.go` to add fresh data
