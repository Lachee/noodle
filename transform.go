package noodle

import "log"

//Transform2D is deprecated
type Transform2D struct {
	Position Vector2
	Rotation float32
	Scale    Vector2
}

func NewTransform2D(position Vector2, rotation float32, scale Vector2) Transform2D {
	return Transform2D{position, rotation, scale}
}

//Transform handles relativeness
type Transform struct {
	parent   *Transform
	children []*Transform

	localPosition Vector3
	localScale    Vector3
	localRotation Quaternion

	//worldMatrix is a cached matrix of the world representation of this object
	worldMatrix     Matrix
	needWorldUpdate bool
}

//NewTransform creates a new blank transform
func NewTransform() *Transform {
	t := &Transform{}
	t.localPosition = Vector3{}
	t.localScale = Vector3{1, 1, 1}
	t.localRotation = NewQuaternionIdentity()
	t.needWorldUpdate = true
	t.children = make([]*Transform, 1)
	t.SetParent(nil)
	return t
}

//GetParent gets the current parent, otherwise nil
func (t *Transform) GetParent() *Transform { return t.parent }

//GetLocalPosition gets the local position
func (t *Transform) GetLocalPosition() Vector3 { return t.localPosition }

//GetLocalScale gets the local scale
func (t *Transform) GetLocalScale() Vector3 { return t.localScale }

//GetLocalRotation gets the local rotation
func (t *Transform) GetLocalRotation() Quaternion { return t.localRotation }

//GetPosition gets the world position
func (t *Transform) GetPosition() Vector3 {
	if t.needWorldUpdate {
		t.updateWorld(false)
	}
	return t.worldMatrix.TransformCoordinate(Vector3{0, 0, 0})
}

//GetWorldMatrix gets teh current world matrix
func (t *Transform) GetWorldMatrix() Matrix {
	if t.needWorldUpdate {
		t.updateWorld(false)
	}
	return t.worldMatrix
}

//GetLocalMatrix returns a local matrix
func (t *Transform) GetLocalMatrix() Matrix {
	translate := NewMatrixTranslate(t.localPosition)
	rotation := NewMatrixRotation(t.localRotation)
	scale := NewMatrixScale(t.localScale)
	result := translate.Multiply(rotation).Multiply(scale)
	log.Println("translate", translate)
	log.Println("rotation", rotation)
	log.Println("scale", scale)
	log.Println("result", result)
	log.Println("================")
	return result
}

//SetParent sets the transform parent.
func (t *Transform) SetParent(parent *Transform) {
	if t.parent != parent {
		if t.parent != nil {
			t.parent.removeChild(t)
			t.parent = nil
		}
		if parent != nil {
			parent.addChild(t)
		}
		t.parent = parent
		t.requireWorldUpdate()
	}
}

//SetLocalPosition sets the local position
func (t *Transform) SetLocalPosition(v Vector3) {
	t.localPosition = v
	t.requireWorldUpdate()
}

//SetLocalScale sets the local scale
func (t *Transform) SetLocalScale(v Vector3) {
	t.localScale = v
	t.requireWorldUpdate()
}

//SetLocalRotation sets the local rotation
func (t *Transform) SetLocalRotation(q Quaternion) {
	t.localRotation = q
	t.requireWorldUpdate()
}

//SetPosition sets the global position
func (t *Transform) SetPosition(v Vector3) {
	if t.parent != nil {
		t.localPosition = v.Subtract(t.parent.GetPosition())
	} else {
		t.localPosition = v
	}
	t.requireWorldUpdate()
}

// Modifiers

//Translate moves the transform
func (t *Transform) Translate(v Vector3) {
	t.localPosition = t.localPosition.Add(v)
	t.requireWorldUpdate()
}

//Rotate rotates the transform
func (t *Transform) Rotate(q Quaternion) {
	t.localRotation = t.localRotation.Multiply(q)
	t.requireWorldUpdate()
}

////LookAt tells the transform to look at a thing
//func (t *Transform) LookAt(v, up Vector3) {
//	v2 := v.Subtract(t.GetPosition())
//	if v2.SqrLength() > 0 {
//		// first build rotation matrix
//		zaxis := v2.Normalize().Negate()
//		xaxis := zaxis.Cross(up).Normalize().Negate()
//		yaxis := zaxis.Cross(xaxis)
//		t.localRotation = NewQuaternionAxis(xaxis, yaxis, zaxis)
//		t.requireWorldUpdate()
//	}
//}

// Utilities

//TransformPoint converts the point from local to world space.
func (t *Transform) TransformPoint(p Vector3) Vector3 {
	return t.GetWorldMatrix().TransformCoordinate(p)
}

func (t *Transform) addChild(child *Transform)    {}
func (t *Transform) removeChild(child *Transform) {}

//updateWorld sets the world matrix
func (t *Transform) updateWorld(includeChildren bool) {

	//Multiply our local matrix by the world matrix. Otherwise its just use our local matrix
	if t.parent != nil {
		t.worldMatrix = t.parent.GetWorldMatrix().Multiply(t.GetLocalMatrix())
	} else {
		t.worldMatrix = t.GetLocalMatrix()
	}

	//We no longer need shit
	t.needWorldUpdate = false

	//Update children if we need too
	if includeChildren {
		for _, child := range t.children {
			child.updateWorld(true)
		}
	}
}

//requireWorldUpdate tells the transform and its children that a update is required on the world axis.
func (t *Transform) requireWorldUpdate() {

	//If its already set, ignore
	if t.needWorldUpdate {
		return
	}

	//Update our state
	t.needWorldUpdate = true

	//Update the children
	for _, child := range t.children {
		if child != nil {
			child.requireWorldUpdate()
		}
	}
}
