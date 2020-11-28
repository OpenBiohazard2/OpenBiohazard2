package render

import (
	"github.com/samuelyuan/openbiohazard2/fileio"
)

const (
	ANIMATION_FRAME_TIME = 30 // time in milliseconds
)

type Animation struct {
	TotalTime   float64
	FrameIndex  int // index in AnimationIndexFrames
	FrameNumber int // corresponds to 1 frame id in the animation loop
	CurPose     int
}

func NewAnimation() *Animation {
	return &Animation{
		TotalTime:   float64(0),
		FrameIndex:  0,
		FrameNumber: 0,
		CurPose:     -1,
	}
}

func (animation *Animation) UpdateAnimationFrame(
	poseNumber int,
	animationData *fileio.EDDOutput,
	timeElapsedSeconds float64,
) {
	// Only keep track of time if an animation is playing
	if animation.CurPose != -1 {
		animation.TotalTime += timeElapsedSeconds * 1000
	} else {
		animation.TotalTime = 0
	}

	// Switch to a different pose
	if animation.CurPose != poseNumber {
		animation.FrameIndex = 0
		if poseNumber != -1 {
			frameData := animationData.AnimationIndexFrames[poseNumber]
			animation.FrameNumber = frameData[animation.FrameIndex].FrameId
		}
		animation.CurPose = poseNumber
	}

	// Loop animation data
	if animation.TotalTime >= ANIMATION_FRAME_TIME && animation.CurPose != -1 {
		animation.TotalTime = 0
		animation.FrameIndex++
		if poseNumber != -1 {
			frameData := animationData.AnimationIndexFrames[poseNumber]
			if animation.FrameIndex >= len(frameData) {
				animation.FrameIndex = 0
			}
			animation.FrameNumber = frameData[animation.FrameIndex].FrameId
		}
	}
}
