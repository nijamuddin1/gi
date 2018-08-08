// Copyright (c) 2018, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package giv

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/goki/gi"
	"github.com/goki/gi/units"
	"github.com/goki/ki"
	"github.com/goki/ki/kit"
)

// MapViewInline represents a map as a single line widget, for smaller maps
// and those explicitly marked inline -- constructs widgets in Parts to show
// the key names and editor vals for each value.
type MapViewInline struct {
	gi.PartsWidgetBase
	Map     interface{} `desc:"the map that we are a view onto"`
	Keys    []ValueView `json:"-" xml:"-" desc:"ValueView representations of the map keys"`
	Values  []ValueView `json:"-" xml:"-" desc:"ValueView representations of the fields"`
	TmpSave ValueView   `json:"-" xml:"-" desc:"value view that needs to have SaveTmp called on it whenever a change is made to one of the underlying values -- pass this down to any sub-views created from a parent"`
	ViewSig ki.Signal   `json:"-" xml:"-" desc:"signal for valueview -- only one signal sent when a value has been set -- all related value views interconnect with each other to update when others update"`
}

var KiT_MapViewInline = kit.Types.AddType(&MapViewInline{}, MapViewInlineProps)

// SetMap sets the source map that we are viewing -- rebuilds the children to represent this map
func (mv *MapViewInline) SetMap(mp interface{}, tmpSave ValueView) {
	// note: because we make new maps, and due to the strangeness of reflect, they
	// end up not being comparable types, so we can't check if equal
	mv.Map = mp
	mv.TmpSave = tmpSave
	mv.UpdateFromMap()
}

var MapViewInlineProps = ki.Props{
	"min-width": units.NewValue(60, units.Ex),
}

// todo: maybe figure out a way to share some of this redundant code..

// ConfigParts configures Parts for the current map
func (mv *MapViewInline) ConfigParts() {
	if kit.IfaceIsNil(mv.Map) {
		return
	}
	mv.Parts.Lay = gi.LayoutHoriz
	config := kit.TypeAndNameList{}
	// always start fresh!
	mv.Keys = make([]ValueView, 0)
	mv.Values = make([]ValueView, 0)

	mpv := reflect.ValueOf(mv.Map)
	mpvnp := kit.NonPtrValue(mpv)

	keys := mpvnp.MapKeys()
	sort.Slice(keys, func(i, j int) bool {
		return kit.ToString(keys[i]) < kit.ToString(keys[j])
	})
	for _, key := range keys {
		kv := ToValueView(key.Interface())
		if kv == nil { // shouldn't happen
			continue
		}
		kv.SetMapKey(key, mv.Map, mv.TmpSave)

		val := mpvnp.MapIndex(key)
		vv := ToValueView(val.Interface())
		if vv == nil { // shouldn't happen
			continue
		}
		vv.SetMapValue(val, mv.Map, key.Interface(), kv, mv.TmpSave) // needs key value view to track updates

		keytxt := kit.ToString(key.Interface())
		keynm := fmt.Sprintf("key-%v", keytxt)
		valnm := fmt.Sprintf("value-%v", keytxt)

		config.Add(kv.WidgetType(), keynm)
		config.Add(vv.WidgetType(), valnm)
		mv.Keys = append(mv.Keys, kv)
		mv.Values = append(mv.Values, vv)
	}
	config.Add(gi.KiT_Action, "edit-action")
	mods, updt := mv.Parts.ConfigChildren(config, false)
	if !mods {
		updt = mv.Parts.UpdateStart()
	}
	for i, vv := range mv.Values {
		vvb := vv.AsValueViewBase()
		vvb.ViewSig.ConnectOnly(mv.This, func(recv, send ki.Ki, sig int64, data interface{}) {
			mvv, _ := recv.EmbeddedStruct(KiT_MapViewInline).(*MapViewInline)
			mvv.ViewSig.Emit(mvv.This, 0, nil)
		})
		keyw := mv.Parts.KnownChild(i * 2).(gi.Node2D)
		widg := mv.Parts.KnownChild((i * 2) + 1).(gi.Node2D)
		kv := mv.Keys[i]
		kv.ConfigWidget(keyw)
		vv.ConfigWidget(widg)
	}
	edack, ok := mv.Parts.Children().ElemFromEnd(0)
	if ok {
		edac := edack.(*gi.Action)
		edac.SetIcon("edit")
		edac.Tooltip = "map edit dialog"
		edac.ActionSig.ConnectOnly(mv.This, func(recv, send ki.Ki, sig int64, data interface{}) {
			mvv, _ := recv.EmbeddedStruct(KiT_MapViewInline).(*MapViewInline)
			dlg := MapViewDialog(mvv.Viewport, mvv.Map, mvv.TmpSave, "Map Value View", "", nil, nil, nil)
			mvvvk, ok := dlg.Frame().Children().ElemByType(KiT_MapView, true, 2)
			if ok {
				mvvv := mvvvk.(*MapView)
				mvvv.ViewSig.ConnectOnly(mvv.This, func(recv, send ki.Ki, sig int64, data interface{}) {
					mvvvv, _ := recv.EmbeddedStruct(KiT_MapViewInline).(*MapViewInline)
					mvvvv.ViewSig.Emit(mvvvv.This, 0, nil)
				})
			}
		})
	}
	mv.Parts.UpdateEnd(updt)
}

func (mv *MapViewInline) UpdateFromMap() {
	mv.ConfigParts()
}

func (mv *MapViewInline) UpdateValues() {
	// maps have to re-read their values because they can't get pointers!
	mv.ConfigParts()
}

func (mv *MapViewInline) Style2D() {
	mv.ConfigParts()
	mv.WidgetBase.Style2D()
}

func (mv *MapViewInline) Render2D() {
	if mv.FullReRenderIfNeeded() {
		return
	}
	if mv.PushBounds() {
		mv.Render2DParts()
		mv.Render2DChildren()
		mv.PopBounds()
	}
}