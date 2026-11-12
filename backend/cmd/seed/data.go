package main

import "math/rand"

var firstNames = []string{
	"Alice", "Bob", "Carol", "David", "Eve", "Frank", "Grace", "Henry",
	"Ivy", "Jack", "Kate", "Liam", "Mia", "Noah", "Olivia", "Peter",
	"Quinn", "Rachel", "Sam", "Tina", "Umar", "Vera", "Will", "Xena",
	"Yara", "Zane", "Adam", "Bella", "Carlos", "Diana", "Elias", "Fiona",
	"Gabriel", "Hana", "Igor", "Julia", "Karl", "Lena", "Marco", "Nora",
	"Oscar", "Paula", "Rafael", "Sofia", "Tomas", "Ulla", "Viktor", "Wendy",
	"Xu", "Yusuf", "Zoe", "Anna", "Ben", "Chloe", "Daniel", "Emma",
	"Felix", "Greta", "Hugo", "Iris", "Jonas", "Klara", "Lucas", "Maya",
	"Niklas", "Olga", "Pavel", "Rita", "Sebastian", "Thea", "Ursula", "Valtteri",
	"Wilma", "Xavier", "Ylva", "Zachary", "Aino", "Eero", "Iida", "Juho",
	"Kaisa", "Lauri", "Mikko", "Noora", "Onni", "Pinja", "Rasmus", "Sanni",
	"Timo", "Venla", "Aleksi", "Emilia", "Joonas", "Krista", "Leo", "Matilda",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller",
	"Davis", "Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez",
	"Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin",
	"Lee", "Perez", "Thompson", "White", "Harris", "Sanchez", "Clark",
	"Ramirez", "Lewis", "Robinson", "Walker", "Young", "Allen", "King",
	"Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores", "Green",
	"Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell", "Mitchell",
	"Carter", "Roberts", "Virtanen", "Korhonen", "Mäkinen", "Nieminen",
	"Hämäläinen", "Laine", "Heikkinen", "Koskinen", "Järvinen", "Lehtonen",
	"Lehtinen", "Saarinen", "Salminen", "Heinonen", "Niemi", "Heikkilä",
	"Kinnunen", "Salonen", "Turunen", "Salo", "Laitinen", "Tuominen",
	"Rantanen", "Karjalainen", "Jokinen", "Mattila", "Savolainen", "Lahtinen",
	"Ahonen", "Ojala", "Leppänen", "Kallio", "Väisänen", "Hiltunen",
	"Manninen", "Koivisto", "Pitkänen", "Räsänen", "Toivonen", "Hakala",
}

var companyPrefixes = []string{
	"Polaris", "Aurora", "Vertex", "Nimbus", "Lumen", "Zenith", "Orbit",
	"Nova", "Atlas", "Quantum", "Cobalt", "Helix", "Nebula", "Fusion",
	"Beacon", "Cascade", "Delta", "Echo", "Forge", "Granite", "Harbor",
	"Ignite", "Juniper", "Keystone", "Lighthouse", "Meridian", "Northstar",
	"Onyx", "Prism", "Quarry", "Ridge", "Summit", "Titan", "Umbra",
	"Vanguard", "Wavelength", "Xenon", "Yonder", "Zenith", "Arctic",
	"Boreal", "Cinder", "Driftwood", "Ember", "Frontier", "Glacier",
	"Highland", "Ironwood", "Journey", "Kestrel",
}

var companySuffixes = []string{
	"Labs", "Systems", "Technologies", "Solutions", "Software", "Studio",
	"Works", "Digital", "Group", "Collective", "Industries", "Dynamics",
	"Analytics", "Networks", "Cloud", "Data", "AI", "Robotics", "Ventures",
	"Partners", "Consulting", "Interactive", "Media", "Financial", "Health",
	"Energy", "Logistics", "Retail", "Foods", "Aerospace",
}

var cities = []string{
	"Helsinki", "Espoo", "Tampere", "Vantaa", "Oulu", "Turku", "Jyväskylä",
	"Lahti", "Kuopio", "Pori", "Stockholm", "Gothenburg", "Malmö", "Oslo",
	"Bergen", "Copenhagen", "Aarhus", "Berlin", "Munich", "Hamburg",
	"Amsterdam", "Rotterdam", "Brussels", "Paris", "Lyon", "Madrid",
	"Barcelona", "Lisbon", "Porto", "Milan", "Rome", "Vienna", "Zurich",
	"Geneva", "Prague", "Warsaw", "Krakow", "Budapest", "Bucharest",
	"Dublin", "London", "Manchester", "Edinburgh", "Toronto", "Vancouver",
	"New York", "San Francisco", "Seattle", "Austin", "Boston",
}

var jobTitles = []string{
	"Software Engineer", "Senior Software Engineer", "Staff Software Engineer",
	"Backend Engineer", "Senior Backend Engineer", "Frontend Engineer",
	"Senior Frontend Engineer", "Full Stack Engineer", "DevOps Engineer",
	"Site Reliability Engineer", "Platform Engineer", "Data Engineer",
	"Senior Data Engineer", "Machine Learning Engineer", "ML Ops Engineer",
	"Mobile Engineer (iOS)", "Mobile Engineer (Android)", "Security Engineer",
	"QA Engineer", "Automation Engineer", "Product Manager",
	"Senior Product Manager", "Technical Product Manager", "Engineering Manager",
	"Director of Engineering", "VP of Engineering", "CTO",
	"Product Designer", "UX Designer", "UI Designer", "Graphic Designer",
	"Data Scientist", "Data Analyst", "Business Analyst", "Solutions Architect",
	"Cloud Architect", "Technical Lead", "Scrum Master", "Project Manager",
	"Program Manager", "Account Manager", "Sales Engineer", "Customer Success Manager",
	"Marketing Manager", "Growth Marketer", "Content Writer", "SEO Specialist",
	"Finance Manager", "HR Business Partner", "Recruiter", "Talent Acquisition Specialist",
}

var jobCategories = []string{
	"Engineering", "Data", "Design", "Product", "Marketing", "Sales",
	"Operations", "Finance", "HR", "Customer Success", "Security", "Infrastructure",
}

var skillTags = []string{
	"golang", "python", "javascript", "typescript", "react", "vue", "svelte",
	"node", "rust", "java", "kotlin", "swift", "ruby", "php", "csharp",
	"postgres", "mysql", "mongodb", "redis", "kafka", "rabbitmq", "elasticsearch",
	"docker", "kubernetes", "terraform", "ansible", "aws", "gcp", "azure",
	"graphql", "rest", "grpc", "microservices", "serverless", "linux",
	"git", "ci-cd", "devops", "sre", "observability", "prometheus", "grafana",
	"machine-learning", "pytorch", "tensorflow", "nlp", "computer-vision",
	"agile", "scrum", "kanban", "tdd", "bdd", "leadership", "mentoring",
	"figma", "sketch", "tailwind", "bootstrap", "sass", "webpack", "vite",
	"sql", "nosql", "etl", "data-warehousing", "spark", "airflow",
}

var companyDescriptions = []string{
	"We build tools that help teams ship faster. Backed by top-tier investors and profitable since 2021.",
	"Engineering-led company focused on developer experience and open-source. Distributed team across Europe.",
	"Series B startup disrupting the logistics space with real-time optimization.",
	"Bootstrapped consultancy with a strong track record of complex integrations. We pick our clients carefully.",
	"Product company serving 2M+ users worldwide. Growing 100% YoY. Small team, big impact.",
	"Enterprise software vendor serving Fortune 500 clients. Stable, profitable, remote-friendly.",
	"Agency working with startups to ship MVPs in weeks. Fast-paced, high autonomy.",
	"Non-profit tech organization building tools for the public sector. Mission-driven team.",
	"Hardware company expanding into IoT. Building the software stack from scratch.",
	"Long-established software house modernizing legacy codebases for enterprise clients.",
}

var seekerHeadlines = []string{
	"Backend developer with 5+ years in Go and Postgres",
	"Full-stack engineer passionate about DX and product craft",
	"Data scientist specializing in NLP and recommendation systems",
	"Frontend engineer focused on accessibility and performance",
	"DevOps engineer automating everything that moves",
	"Mobile engineer with shipped apps in the App Store top 100",
	"Product manager bridging user research and engineering",
	"Engineering manager building high-trust teams",
	"Recent CS graduate looking for my first role",
	"Self-taught developer with a portfolio of side projects",
	"Cloud architect specializing in AWS and Kubernetes",
	"Security researcher with a focus on application security",
}

var seekerBios = []string{
	"I've spent the last few years building products end-to-end. I care about clean code, good tests, and shipping.\n\nIn my spare time I contribute to open source and write about engineering practices on my blog.",
	"Previously at two early-stage startups. Comfortable wearing many hats and working with ambiguity.\n\nLooking for a team that values autonomy and craft over process.",
	"Transitioning from academia to industry. Strong fundamentals in algorithms and systems.\n\nInterested in roles where I can keep learning while shipping real products.",
	"Senior engineer who's mentored dozens of juniors. I believe good code review is a superpower.\n\nRemote-first for the last four years.",
	"Product-minded engineer. I've shipped features that moved the needle on real business metrics.\n\nNot afraid to write docs, talk to users, or clean up a legacy codebase.",
}

var searchNames = []string{
	"Remote Go jobs", "Frontend in Helsinki", "Senior backend roles",
	"Data science jobs", "Product roles at startups", "Remote-first companies",
	"Junior-friendly positions", "Kubernetes and cloud jobs", "Design roles in Europe",
	"Engineering manager openings",
}

var coverLetters = []string{
	"I'm excited about this role. I've been working with the technologies mentioned in the listing for years and I think I could contribute from day one.",
	"I've been following your company for a while and I love what you're building. I'd be thrilled to be part of the team.",
	"Briefly: I'm looking for a role where I can grow, ship, and work with smart people. Yours looks like a great fit.",
	"Attached is my resume. I'm happy to jump on a call to discuss how I can help.",
	"I applied because the problem space genuinely interests me. Here's a short summary of relevant work I've done.",
}

var jobDescriptions = []string{
	"We're hiring a talented engineer to join our growing team. You'll work on production systems serving millions of requests per day, collaborate with product and design, and help shape the technical direction of the team.\n\nWhat you'll do:\n- Design, build, and ship features end-to-end\n- Write well-tested, maintainable code\n- Participate in code review and technical discussions\n- Mentor more junior engineers\n\nWe value curiosity, clarity in communication, and a bias toward shipping.",
	"As part of our platform team, you'll work on the infrastructure that powers our products. You'll collaborate closely with application engineers to build reliable, scalable systems.\n\nResponsibilities:\n- Own services from design through production\n- Improve observability and reliability\n- Automate repetitive tasks\n- Help define engineering standards\n\nWe're remote-first and support flexible hours.",
	"We're looking for an experienced engineer to help us scale. This is a hands-on role with significant ownership — you'll be designing systems, not just implementing tickets.\n\nYou'll work with:\n- A small, senior team\n- Modern stack (no legacy)\n- A product that users love\n\nWe ship weekly, we don't do death marches, and we invest in tooling.",
	"Join a mission-driven team building tools for an underserved market. We care about craft, we care about users, and we care about each other.\n\nWhat we offer:\n- Competitive salary and equity\n- Flexible work hours and location\n- Learning budget\n- A team that ships\n\nWhat we ask:\n- Strong fundamentals\n- Willingness to learn\n- Clear communication",
}

var jobRequirements = []string{
	"- 3+ years of professional experience\n- Strong fundamentals in computer science\n- Experience with production systems\n- Good written and verbal communication in English",
	"- 5+ years of experience\n- Deep knowledge of at least one language\n- Experience designing and operating distributed systems\n- Track record of shipping complex features",
	"- Solid understanding of web fundamentals\n- Experience with a modern frontend framework\n- Eye for design and attention to detail\n- Comfort working with ambiguity",
	"- Bachelor's degree or equivalent experience\n- Strong problem-solving skills\n- Eagerness to learn\n- Team-oriented mindset",
}

var employmentTypes = []string{"full_time", "full_time", "full_time", "full_time", "part_time", "contract", "internship"}
var remoteTypes = []string{"on_site", "hybrid", "remote", "remote", "remote"}
var currencies = []string{"EUR", "EUR", "EUR", "USD", "GBP"}

func pick(r *rand.Rand, pool []string) string {
	return pool[r.Intn(len(pool))]
}

func pickInt(r *rand.Rand, min, max int) int {
	if max <= min {
		return min
	}
	return min + r.Intn(max-min+1)
}

func pickTags(r *rand.Rand, min, max int) []string {
	n := pickInt(r, min, max)
	seen := make(map[string]bool, n)
	out := make([]string, 0, n)
	for len(out) < n {
		t := pick(r, skillTags)
		if seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}

func randomFullName(r *rand.Rand) string {
	return pick(r, firstNames) + " " + pick(r, lastNames)
}

func randomUsername(r *rand.Rand, i int) string {
	return pick(r, firstNames) + "_" + pick(r, lastNames) + "_" + itoa(i)
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	neg := i < 0
	if neg {
		i = -i
	}
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
