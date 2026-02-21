package event

import (
	"github.com/gofiber/fiber/v3"

	"backend/models"
	"backend/repository"
	"backend/utils"
)

var artistRepo = repository.NewArtistRepository()

func GetAllArtists(c fiber.Ctx) error {
	artists, err := artistRepo.GetAll(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to fetch artists")
	}

	return utils.SuccessResponse(c, 200, "Artists fetched successfully", artists)
}

func GetArtistByID(c fiber.Ctx) error {
	id := c.Params("id")
	artist, err := artistRepo.FindByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, 404, "Artist not found")
	}
	return utils.SuccessResponse(c, 200, "Artist fetched successfully", artist)
}

func CreateArtist(c fiber.Ctx) error {
	var artist models.Artist
	if err := c.Bind().Body(&artist); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if err := artistRepo.Create(c.Context(), &artist); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create artist")
	}

	return utils.SuccessResponse(c, 201, "Artist created successfully", artist)
}

func SeedArtists(c fiber.Ctx) error {
	artists := []models.Artist{
		{
			Name: "DJ Aurora", Role: "Main Artist", ImageURL: "/events/artists/1.png",
			Description: "Electronic music pioneer with over a decade of experience in the underground scene.",
			Genre:       "Electronic, House, Techno",
			IsVerified:  true, Rating: 4.8, ReviewCount: 124, EventsHosted: 45,
			Location: "Mumbai, Maharashtra", ContactEmail: "aurora@dj.com", FollowerCount: 15400,
			ExperienceYears: 12, Specialties: []string{"Deep House", "Melodic Techno"},
			SocialLinks: []string{"https://instagram.com/djaurora", "https://facebook.com/djaurora"},
		},
		{
			Name: "The Waves", Role: "Band", ImageURL: "/events/artists/2.png",
			Description: "Indie rock sensation known for their high-energy performances and soulful lyrics.",
			Genre:       "Indie Rock, Alternative",
			IsVerified:  true, Rating: 4.9, ReviewCount: 89, EventsHosted: 22,
			Location: "Bangalore, Karnataka", ContactEmail: "contact@thewaves.band", FollowerCount: 8200,
			ExperienceYears: 5, Specialties: []string{"Live Band", "Acoustic Sets"},
			SocialLinks: []string{"https://instagram.com/thewavesband", "https://youtube.com/thewaves"},
		},
		{
			Name: "Luna Smith", Role: "Solo Singer", ImageURL: "/events/artists/3.png",
			Description: "Vibrant pop vocalist with a focus on contemporary sounds and empowering themes.",
			Genre:       "Pop, Soul",
			IsVerified:  false, Rating: 4.5, ReviewCount: 56, EventsHosted: 15,
			Location: "Chennai, Tamil Nadu", ContactEmail: "luna@smith.com", FollowerCount: 3100,
			ExperienceYears: 3, Specialties: []string{"Pop Vocals", "Soulful Covers"},
		},
		{
			Name: "Echo Band", Role: "Band", ImageURL: "/events/artists/4.png",
			Description: "Experimental jazz ensemble pushing the boundaries of traditional rhythms.",
			Genre:       "Jazz, Fusion",
			IsVerified:  true, Rating: 4.7, ReviewCount: 42, EventsHosted: 18,
			Location: "Kolkata, West Bengal", ContactEmail: "info@echoband.com", FollowerCount: 2500,
			ExperienceYears: 8, Specialties: []string{"Progressive Jazz", "Fusion"},
		},
		{
			Name: "Sarah Jazz", Role: "Jazz Singer", ImageURL: "/events/artists/5.png",
			Description: "Classic jazz vocalist delivering timeless renditions of legendary standards.",
			Genre:       "Jazz, Blues",
			IsVerified:  true, Rating: 4.9, ReviewCount: 210, EventsHosted: 60,
			Location: "Delhi", ContactEmail: "sarah@jazz.com", FollowerCount: 12000,
			ExperienceYears: 15, Specialties: []string{"Vocal Jazz", "Classic Standards"},
		},
		{
			Name: "Max Power", Role: "Rock Artist", ImageURL: "/events/artists/6.png",
			Description: "High-octane rock performer dedicated to keeping the spirit of rock and roll alive.",
			Genre:       "Rock, Hard Rock",
			IsVerified:  false, Rating: 4.2, ReviewCount: 34, EventsHosted: 10,
			Location: "Pune, Maharashtra", ContactEmail: "max@power.rock", FollowerCount: 1800,
			ExperienceYears: 4, Specialties: []string{"Hard Rock", "Metal"},
		},
	}

	for _, a := range artists {
		_ = artistRepo.Create(c.Context(), &a)
	}

	return utils.SuccessResponse(c, 200, "Artists seeded successfully", nil)
}
