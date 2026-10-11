package checker

func (links *ValueSymbolLinks) getWriteType() *Type {
	if links.rare == nil {
		return nil
	}
	return links.rare.writeType
}

func (links *ValueSymbolLinks) setWriteType(value *Type) {
	if links.rare == nil {
		if value == nil {
			return
		}
		links.rare = &valueSymbolLinksRare{}
	}
	links.rare.writeType = value
}

func (links *ValueSymbolLinks) getNameType() *Type {
	if links.rare == nil {
		return nil
	}
	return links.rare.nameType
}

func (links *ValueSymbolLinks) setNameType(value *Type) {
	if links.rare == nil {
		if value == nil {
			return
		}
		links.rare = &valueSymbolLinksRare{}
	}
	links.rare.nameType = value
}

func (links *ValueSymbolLinks) getContainingType() *Type {
	if links.rare == nil {
		return nil
	}
	return links.rare.containingType
}

func (links *ValueSymbolLinks) setContainingType(value *Type) {
	if links.rare == nil {
		if value == nil {
			return
		}
		links.rare = &valueSymbolLinksRare{}
	}
	links.rare.containingType = value
}

func (links *ValueSymbolLinks) getFunctionOrConstructorChecked() bool {
	if links.rare == nil {
		return false
	}
	return links.rare.functionOrConstructorChecked
}

func (links *ValueSymbolLinks) setFunctionOrConstructorChecked(value bool) {
	if links.rare == nil {
		if value == false {
			return
		}
		links.rare = &valueSymbolLinksRare{}
	}
	links.rare.functionOrConstructorChecked = value
}
