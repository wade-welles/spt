package spt

import (
	"log"
	"os"
	"runtime/pprof"
	"sync"
	"testing"
	"time"
)

func parabolicBowl(radius, height, baseThickness float64) SDF3 {
	shape := Revolve(0, Parabola(radius*2, height))
	return RotateX(270, Difference(shape, TranslateY(baseThickness, shape)))
}

func testScene() Scene {

	steel := Metal(Color{0.4, 0.4, 0.4}, 0.95)
	stainless := Metal(Color{0.4, 0.4, 0.4}, 0.3)
	gold := Metal(Color{0.93, 0.78, 0.31}, 0.0)
	copper := Metal(Color{0.68, 0.45, 0.41}, 0.8)
	brass := Metal(Color{0.80, 0.58, 0.45}, 0.9)

	things := func(material Material, origin Vec3) []Thing {

		at := func(v Vec3, s SDF3) SDF3 {
			return Translate(origin.Add(v), s)
		}

		return append(
			[]Thing{
				Object(material,
					at(V3(0, 0, 500), Sphere(500)),
				),
				Object(material,
					at(V3(0, 2000, 500), Cube(1000, 1000, 1000)),
				),
				Object(material,
					at(V3(-1500, 0, 500), Cylinder(1000, 500)),
				),
				Object(material,
					at(V3(1500, 2000, 500), Cone(1000, 500)),
				),
				Object(material,
					at(V3(1500, 0, 500), Torus(500, 350)),
				),
				Object(material,
					at(V3(-1500, 2000, 500), Pyramid(1000, 1000)),
				),
			}, func() (set []Thing) {
				for i := 3; i <= 8; i++ {
					set = append(set, Object(material,
						at(V3(float64(i-3)*600-1450, -1250, 250), Extrude(500, Polygon(i, 250))),
					))
				}
				return
			}()...,
		)
	}

	stuff := []Thing{
		WorkBench(25000),
		Object(
			Light(White.Scale(4)),
			Translate(V3(-7500, 0, 20000), Sphere(10000)),
		),
		Object(
			brass,
			Translate(V3(-5500, -1200, 500),
				Difference(
					Intersection(
						Sphere(500),
						Cube(800, 800, 800),
					),
					Cylinder(1002, 200),
					RotateX(90, Cylinder(1002, 200)),
					RotateY(90, Cylinder(1002, 200)),
				),
			),
		),
		Object(
			brass,
			Translate(V3(5500, -1200, 500),
				Union(
					Intersection(
						Sphere(500),
						Cube(800, 800, 800),
					),
					Cylinder(1250, 200),
					RotateX(90, Cylinder(1250, 200)),
					RotateY(90, Cylinder(1250, 200)),
				),
			),
		),
		Object(
			brass,
			Translate(V3(6500, 3750, 500),
				Round(100, Cube(800, 800, 800)),
			),
		),
		Object(
			brass,
			Translate(V3(-6500, 3750, 500),
				Round(100, Cylinder(800, 500)),
			),
		),
		Object(
			steel,
			Translate(V3(6000, 900, 600),
				Repeat(1, 1, 1, 400, 400, 400, Sphere(200)),
			),
		),
		Object(
			copper,
			Translate(V3(-6000, 900, 750),
				Capsule(750, 500, 250),
			),
		),
		Object(
			Glass(Color{1.0, 0.5, 0.5}, 1.5),
			Union(
				TranslateZ(100, CylinderR(200, 500, 10)),
				parabolicBowl(1000, 2000, 200),
			),
		),
	}

	stuff = append(stuff, things(steel, V3(-2750, -2750, 0))...)
	stuff = append(stuff, things(copper, V3(2750, -2750, 0))...)
	stuff = append(stuff, things(stainless, V3(-3000, 2250, 0))...)
	stuff = append(stuff, things(gold, V3(3000, 2250, 0))...)

	for i := 0; i < 9; i++ {
		color := Color{0.5, 0.5, 1.0}
		if i%2 == 0 {
			color = Color{0.5, 1.0, 0.5}
		}
		if i >= 3 && i < 6 {
			continue
		}
		stuff = append(stuff, Object(
			Glass(color, 1.5),
			Translate(V3(0, float64(i)*1000-3500, 250), Sphere(250)),
		))
	}

	return Scene{
		Width:     1280,
		Height:    720,
		Passes:    10,
		Samples:   1,
		Bounces:   16,
		Horizon:   100000,
		Threshold: 0.0001,
		Ambient:   White.Scale(0.05),

		Camera: NewCamera(
			V3(0, -8000, 8000),
			V3(0, -1000, 500),
			Z3,
			40,
			Zero3,
			0.0,
		),

		Stuff: stuff,
	}
}

func TestLocal(t *testing.T) {

	cf, cerr := os.Create("cpuprofile")
	if cerr != nil {
		log.Fatal(cerr)
	}
	pprof.StartCPUProfile(cf)
	defer pprof.StopCPUProfile()

	hf, herr := os.Create("heapprofile")
	if herr != nil {
		log.Fatal(herr)
	}
	defer pprof.WriteHeapProfile(hf)

	RenderSave("test.png", testScene(), []Renderer{NewLocalRenderer()})
}

func TestRPC(t *testing.T) {
	var group sync.WaitGroup
	stop := make(chan struct{})

	group.Add(1)
	go func() {
		RenderServeRPC(stop, 34242)
		group.Done()
	}()

	time.Sleep(time.Second)
	RenderSave("test.png", testScene(), []Renderer{
		NewRPCRenderer("127.0.0.1:34242"),
	})

	close(stop)
	group.Wait()
}

func TestRPC2(t *testing.T) {
	scene := testScene()
	scene.Width = 1920
	scene.Height = 1080
	scene.Passes = 0
	RenderSave("test.png", scene, []Renderer{
		NewRPCRenderer("slave1:34242"),
		NewRPCRenderer("slave2:34242"),
	})
}

func TestShadow(t *testing.T) {
	scene := Scene{
		Width:     1280,
		Height:    720,
		Passes:    1000,
		Samples:   1,
		Bounces:   8,
		Horizon:   100000,
		Threshold: 0.0001,
		Ambient:   White.Scale(0.05),
		ShadowD:   2.0,
		ShadowR:   0.5,
		ShadowH:   0.8,
		ShadowL:   0.2,

		Camera: NewCamera(
			V3(0, -8000, 8000),
			V3(0, -1000, 500),
			Z3,
			40,
			Zero3,
			0.0,
		),

		Stuff: []Thing{
			SpaceTime(25000),
			//WorkBench(25000),
			Object(
				Light(White.Scale(4)),
				Translate(V3(-7500, 0, 20000), Sphere(10000)),
			),
			Object(
				Matt(White),
				Translate(V3(0, 0, 1000), Sphere(1000)),
			),
		},
	}

	RenderSave("test.png", scene, []Renderer{NewLocalRenderer()})
}

func TestSDF(t *testing.T) {

	scene := Scene{
		Width:     1280,
		Height:    720,
		Passes:    10,
		Samples:   1,
		Bounces:   8,
		Horizon:   100000,
		Threshold: 0.0001,
		Ambient:   White.Scale(0.05),

		Camera: NewCamera(
			V3(0, -8000, 8000),
			V3(0, 0, 500),
			Z3,
			40,
			Zero3,
			0.0,
		),

		Stuff: []Thing{
			WorkBench(50000),
			Object(
				Light(White.Scale(4)),
				Translate(V3(-7500, 0, 20000), Sphere(10000)),
			),
			Object(
				Copper,
				TranslateZ(500, Elongate(500, 0, 0, Cylinder(1000, 500))),
			),
		},
	}

	RenderSave("test.png", scene, []Renderer{NewLocalRenderer()})
}
