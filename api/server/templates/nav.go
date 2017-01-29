package templates

type BasePage struct {
	NavSections []*NavSection
	NavHeader   bool
	MiniNavBar  bool
	Title       string
	Body        Body
}

func (b *BasePage) GetNavSection(name string) *NavSection {
	for _, sec := range b.NavSections {
		if sec.Name == name {
			return sec
		}
	}
	return nil
}
func (b *BasePage) NavSection(name string) *NavSection {
	sec := b.GetNavSection(name)
	if sec != nil {
		return sec
	}
	return b.AddSection(name)
}

func (b *BasePage) AddSection(name string) *NavSection {
	sec := &NavSection{
		Title: name,
		Name:  name,
		Href:  "#",
		Icon:  "th-large",
	}
	b.NavSections = append(b.NavSections, sec)
	return sec
}

type NavSection struct {
	Name     string
	Title    string
	Href     string
	Icon     string
	Active   bool
	NavLinks []*NavLink
}

func (s *NavSection) SetTitle(title string) *NavSection {
	s.Title = title
	return s
}
func (s *NavSection) SetHref(href string) *NavSection {
	s.Href = href
	return s
}
func (s *NavSection) SetActive(active bool) *NavSection {
	s.Active = active
	return s
}

func (s *NavSection) SetIcon(icon string) *NavSection {
	s.Icon = icon
	return s
}

func (s *NavSection) GetNavLink(name string) *NavLink {
	for _, link := range s.NavLinks {
		if link.Name == name {
			return link
		}
	}
	return nil
}

func (s *NavSection) NavLink(name string) *NavLink {
	link := s.GetNavLink(name)
	if link != nil {
		return link
	}
	return s.AddLink(name)
}

func (s *NavSection) AddLink(name string) *NavLink {
	link := &NavLink{Name: name, Title: name, Href: "#"}
	s.NavLinks = append(s.NavLinks, link)
	return link
}

type NavLink struct {
	Title  string
	Name   string
	Href   string
	Active bool
}

func (l *NavLink) SetTitle(title string) *NavLink {
	l.Title = title
	return l
}

func (l *NavLink) SetActive(active bool) *NavLink {
	l.Active = active
	return l
}

func (l *NavLink) SetHref(href string) *NavLink {
	l.Href = href
	return l
}
