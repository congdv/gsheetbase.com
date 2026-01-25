# Google Sheets CMS Templates

## 1. Link-in-Bio

  Field         Description
  ------------- ------------------------
  id            Unique identifier
  title         Link title
  url           Destination URL
  icon          Icon name or emoji
  order         Display order
  active        true/false
  description   Optional subtitle text

## 2. SaaS / Startup Landing Page

  Field           Description
  --------------- -----------------------
  id              Unique section id
  hero_title      Main headline
  hero_subtitle   Supporting text
  cta_text        Call-to-action label
  cta_url         CTA destination
  feature_title   Feature name
  feature_desc    Feature description
  testimonial     User testimonial text
  author_name     Testimonial author
  pricing_plan    Plan name
  price           Plan price

## 3. Waitlist / Pre-launch Page

  Field           Description
  --------------- ----------------------------
  id              Unique id
  product_name    Product name
  tagline         Short pitch
  value_prop      Key benefit
  email           Subscriber email
  status          invited / pending / joined
  referral_code   Optional referral
  launch_date     Planned launch date

## 4. Personal Portfolio / CV

  Field          Description
  -------------- ---------------------
  id             Unique id
  name           Person name
  title          Job title
  bio            Short bio
  skill          Skill name
  project_name   Project title
  project_desc   Project description
  project_url    Link to project
  company        Company name
  role           Role at company
  start_date     Start date
  end_date       End date

## 5. Services / Freelancer Page

  Field          Description
  -------------- ---------------------
  id             Unique id
  service_name   Name of service
  description    Service description
  price          Service price
  billing_type   one-time / monthly
  package_name   Package tier
  turnaround     Delivery time
  testimonial    Client feedback
  client_name    Client name
  booking_url    Booking link

## 6. Directory / List Site

  Field         Description
  ------------- ----------------------
  id            Unique id
  name          Item name
  description   Short description
  category      Category tag
  website_url   External link
  logo_url      Logo image
  featured      true/false
  tags          Comma separated tags
  created_at    Date added

## 7. Job Board

  Field         Description
  ------------- -------------------------------
  id            Job id
  job_title     Position title
  company       Company name
  location      Job location
  salary        Salary range
  job_type      full-time / contract / remote
  description   Job description
  apply_url     Application link
  status        open / closed
  posted_date   Date posted

## 8. Event / Conference Page

  Field          Description
  -------------- --------------------
  id             Unique id
  event_name     Event name
  date           Event date
  time           Event time
  location       Physical or online
  speaker_name   Speaker name
  speaker_bio    Speaker bio
  agenda_item    Schedule item
  ticket_url     Ticket link
  sponsor_name   Sponsor

## 9. Testimonials Wall

  Field         Description
  ------------- ------------------
  id            Unique id
  name          Person name
  role          Job title
  company       Company
  testimonial   Testimonial text
  rating        1--5
  avatar_url    Profile image
  featured      true/false

## 10. Newsletter Landing Page

  Field             Description
  ----------------- --------------------
  id                Unique id
  newsletter_name   Name of newsletter
  headline          Main pitch
  description       Long description
  issue_title       Past issue title
  issue_url         Past issue link
  subscribers       Subscriber count
  signup_url        Email signup link

## 11. Single Product Store

  Field          Description
  -------------- ----------------------
  id             SKU
  product_name   Product name
  description    Product description
  price          Price
  currency       Currency
  image_url      Product image
  stripe_url     Stripe checkout link
  stock          Stock quantity
  category       Product category

## 12. Real Estate Listing

  Field         Description
  ------------- ---------------------
  id            Listing id
  address       Property address
  price         Listing price
  bedrooms      Number of bedrooms
  bathrooms     Number of bathrooms
  area_sqft     Size
  description   Description
  image_url     Property image
  agent_name    Agent
  contact_url   Contact link

## 13. Content / Blog (Headless)

  Field          Description
  -------------- ---------------
  id             Post id
  title          Post title
  slug           URL slug
  excerpt        Short summary
  content        Full content
  author         Author name
  published_at   Publish date
  cover_image    Cover image
  tags           Tags

## 14. Pricing Table

  Field       Description
  ----------- ------------------
  id          Plan id
  plan_name   Plan name
  price       Price
  billing     monthly / yearly
  feature     Feature item
  highlight   true/false
  cta_text    Button label
  cta_url     Checkout link

## 15. FAQ Page

  Field      Description
  ---------- ---------------
  id         Unique id
  question   FAQ question
  answer     FAQ answer
  category   Grouping
  order      Display order
  visible    true/false

## 16. Product Catalog (E-commerce)

| Field         | Description                       |
|---------------|-----------------------------------|
| id            | Unique identifier/SKU             |
| title         | Product name                      |
| description   | Marketing copy                    |
| price         | Numeric price                     |
| currency      | e.g., USD, EUR                    |
| image_url     | Link to product photo             |
| product_url   | Direct link to buy                |
| stock_status  | e.g., "In Stock", "Out of Stock" |
| category      | e.g., "Apparel", "Home"           |

---

## 17. Lead Gen / Form Capture

| Field         | Description                        |
|---------------|------------------------------------|
| timestamp     | Date and time of entry             |
| full_name     | User's name                        |
| email         | User's contact address             |
| interest_level| e.g., "High", "Low"                |
| source_url    | Which page they signed up on       |
| status        | e.g., "New", "Replied", "Closed"   |
| notes         | User's message or internal comments|

---

## 18. Resource Masterlist (Link-in-Bio)

| Field         | Description                        |
|---------------|------------------------------------|
| name          | Resource title                     |
| tagline       | Short description                  |
| link          | Destination URL                    |
| category      | e.g., "Tools", "Articles"          |
| is_featured   | Boolean (TRUE/FALSE)               |
| button_text   | e.g., "Get Started"                |
| icon_code     | e.g., FontAwesome class name       |

---

## 19. Service Menu (Salons & Spas)

| Field         | Description                                 |
|---------------|---------------------------------------------|
| category      | e.g., "Manicures", "Facials"                |
| service_name  | Specific treatment name                     |
| price         | Cost of service                             |
| duration_mins | Time required (e.g., 30, 60)                |
| description   | What’s included                             |
| is_popular    | Boolean (adds a "Best Seller" badge)        |
| booking_url   | Direct link to book                         |

---

## 20. Restaurant Menu

| Field         | Description                        |
|---------------|------------------------------------|
| section       | e.g., "Appetizers", "Drinks"        |
| item_name     | Name of dish/drink                 |
| price         | Numeric price                      |
| ingredients   | List of components                 |
| dietary_tags  | e.g., "Vegan", "Gluten-Free"        |
| spice_level   | Scale of 0-3                       |
| photo_url     | High-quality food image link        |
| available     | Boolean (to hide items)             |

---

## 21. Trade Service Catalog (Plumbers/HVAC)

| Field               | Description                        |
|---------------------|------------------------------------|
| service_id          | Internal code                      |
| category            | e.g., "Emergency", "Installation"  |
| service_name        | Name of the repair/service         |
| base_price          | Flat rate or starting price        |
| emergency_surcharge | Extra fee for after-hours          |
| unit                | e.g., "Per Hour", "Per Visit"      |
| is_emergency        | Boolean                            |
| license_required    | e.g., "Master Plumber Required"    |

---

## 22. Contractor Project Portfolio

| Field            | Description                |
|------------------|---------------------------|
| project_title    | Name of the project        |
| location         | City or neighborhood       |
| service_type     | e.g., "Kitchen Remodel"    |
| completion_date  | Date finished              |
| before_image_url | Link to "Before" photo     |
| after_image_url  | Link to "After" photo      |
| client_quote     | Testimonial text           |
| budget_range     | e.g., "$5,000 - $10,000"   |



# IDEAS:
Phase 1 (Your MVP)

1. Link-in-Bio

2. Directory

3. Job Board

4. Lead Gen

5. Services Page

Phase 2 (Expansion)

6. Single Product Store

7. Newsletter

8. Salon Menu

9. Trade Services


Phase 3 (Bundled Features)

10. Testimonials

11. Pricing

12. FAQ