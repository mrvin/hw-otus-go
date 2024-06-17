package hw03frequencyanalysis

import (
	"reflect"
	"sort"
	"testing"
)

// Change to true if needed
var taskWithAsteriskIsCompleted = true

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func testTop10Subset(t *testing.T, text string, want []string) {
	got := Top10(text)

	if numTOP < len(want) {
		for _, word := range got {
			for i, wordWant := range want {
				if word == wordWant {
					break
				}
				if i == len(want)-1 {
					t.Errorf("got: %v; want: %v", got, want)
				}
			}
		}
	} else {
		t.Errorf("%d < %d numTop >= len(want)", numTOP, len(want))
	}
}

func testTop10Match(t *testing.T, text string, want []string) {
	got := Top10(text)

	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v(after sorting); want: %v(after sorting)", got, want)
	}
}

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		testTop10Match(t, "", []string{})
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{"он", "а", "и", "что", "ты", "не", "если", "то", "его", "кристофер", "робин", "в"}
			testTop10Subset(t, text, expected)
		} else {
			expected := []string{"он", "и", "а", "что", "ты", "не", "если", "-", "то", "Кристофер"}
			testTop10Match(t, text, expected)
		}
	})

	t.Run("10 words", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			text := `All work and no play makes Jack a dull boy.
					All work and no play makes Jack a dull boy.
					All work and no play makes Jack a dull boy.`
			expected := []string{"all", "work", "and", "no", "play", "makes", "jack", "a", "dull", "boy"}

			testTop10Match(t, text, expected)
		}
	})

	t.Run("1 words", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			text := `All.`
			expected := []string{"all"}

			testTop10Match(t, text, expected)
		}
	})

	t.Run("Different forms of the word", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			text := `Нога нога нога! 'нога' (нога) -нога- "нога"`
			expected := []string{"нога"}

			testTop10Match(t, text, expected)
		}
	})

	t.Run("word starts with a dash", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			text := `Но увы -сосредоточиться-то`
			expected := []string{"но", "увы", "сосредоточиться-то"}

			testTop10Match(t, text, expected)
		}
	})

	t.Run("word ends with a dash", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			text := `Но увы- сосредоточиться-то`
			expected := []string{"но", "увы", "сосредоточиться-то"}

			testTop10Match(t, text, expected)
		}
	})

}
