package thinkdatatype

type ThinkGraphNode struct {
	Data any               `json:"-"`
	Edge []*ThinkGraphEdge `json:"-"`
}

type ThinkGraphEdge struct {
	Src      *ThinkGraphNode `json:"-"`
	Dest     *ThinkGraphNode `json:"-"`
	Distance int64           `json:"-"`
}

type ThinkGraph struct {
	m_lstNode []*ThinkGraphNode
}

func NewThinkGraph() *ThinkGraph {
	return &ThinkGraph{
		m_lstNode: make([]*ThinkGraphNode, 0),
	}
}

func (this *ThinkGraph) AddNode(pData any) {
	pNode := &ThinkGraphNode{
		Data: pData,
		Edge: make([]*ThinkGraphEdge, 0),
	}

	this.m_lstNode = append(this.m_lstNode, pNode)
}

func (this *ThinkGraph) edgeExists(pSrc *ThinkGraphNode, pDest *ThinkGraphNode) *ThinkGraphEdge {
	if nil == pSrc || nil == pDest {
		return nil
	}

	for _, pEdge := range pSrc.Edge {
		if pDest == pEdge.Dest {
			return pEdge
		}
	}

	return nil
}

func (this *ThinkGraph) addEdge(pSrc *ThinkGraphNode, pDest *ThinkGraphNode, nDistance int64) {
	if nil == pSrc || nil == pDest || nDistance <= 0 {
		return
	}

	if pEdge := this.edgeExists(pSrc, pDest); pEdge != nil {
		pEdge.Distance = nDistance
		return
	}

	pEdge := &ThinkGraphEdge{
		Src:      pSrc,
		Dest:     pDest,
		Distance: nDistance,
	}

	pSrc.Edge = append(pSrc.Edge, pEdge)
}

func (this *ThinkGraph) nodeExists(pData any) *ThinkGraphNode {
	for _, pNode := range this.m_lstNode {
		if pNode.Data == pData {
			return pNode
		}
	}

	return nil
}

func (this *ThinkGraph) AddEdge(pSrc any, pDest any, nDistance int64) {
	pSrcNode := this.nodeExists(pSrc)
	if nil == pSrcNode {
		pSrcNode = &ThinkGraphNode{
			Data: pSrc,
			Edge: make([]*ThinkGraphEdge, 0),
		}
	}

	pDestNode := this.nodeExists(pDest)
	if nil == pDestNode {
		pDestNode = &ThinkGraphNode{
			Data: pDest,
			Edge: make([]*ThinkGraphEdge, 0),
		}
	}

	this.addEdge(pSrcNode, pDestNode, nDistance)
}



