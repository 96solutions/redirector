package interactor

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lroman242/redirector/domain/dto"
	"github.com/lroman242/redirector/domain/entity"
	"github.com/lroman242/redirector/mocks"
	"net"
	"reflect"
	"testing"
)

func TestClickHandlerFunc_HandleClick(t *testing.T) {
	expectedClick := &entity.Click{
		ID:          "expectedClickID",
		TargetURL:   "expectedTargetURL",
		TRKLink:     nil,
		UserAgent:   nil,
		CountryCode: "expectedCountryCode",
		IP:          net.ParseIP("123.1.2.3"),
		Referer:     "expectedReferer",
		P1:          "expectedP1",
		P2:          "expectedP2",
		P3:          "expectedP3",
		P4:          "expectedP4",
	}

	expFunc := func(ctx context.Context, click *entity.Click) <-chan *dto.ClickProcessingResult {
		ch := make(chan *dto.ClickProcessingResult)
		defer close(ch)

		if expectedClick.ID != click.ID {
			t.Errorf("unexpected click id. expected %s but got %s\n", expectedClick.ID, click.ID)
		}
		if expectedClick.TargetURL != click.TargetURL {
			t.Errorf("unexpected click target URL. expected %s but got %s\n", expectedClick.TargetURL, click.TargetURL)
		}
		if expectedClick.IP.String() != click.IP.String() {
			t.Errorf("unexpected click ip. expected %s but got %s\n", expectedClick.IP.String(), click.IP.String())
		}
		if !reflect.DeepEqual(expectedClick.UserAgent, click.UserAgent) {
			t.Errorf("unexpected click user agent. expected %+v but got %+v\n", expectedClick.UserAgent, click.UserAgent)
		}
		if !reflect.DeepEqual(expectedClick.TRKLink, click.TRKLink) {
			t.Errorf("unexpected click tracking link. expected %+v but got %+v\n", expectedClick.TRKLink, click.TRKLink)
		}
		if expectedClick.CountryCode != click.CountryCode {
			t.Errorf("unexpected click country code. expected %s but got %s\n", expectedClick.CountryCode, click.CountryCode)
		}
		if expectedClick.Referer != click.Referer {
			t.Errorf("unexpected click referer. expected %s but got %s\n", expectedClick.Referer, click.Referer)
		}
		if expectedClick.P1 != click.P1 {
			t.Errorf("unexpected click p1. expected %s but got %s\n", expectedClick.P1, click.P1)
		}
		if expectedClick.P2 != click.P2 {
			t.Errorf("unexpected click p2. expected %s but got %s\n", expectedClick.P2, click.P2)
		}
		if expectedClick.P3 != click.P3 {
			t.Errorf("unexpected click p3. expected %s but got %s\n", expectedClick.P3, click.P3)
		}
		if expectedClick.P4 != click.P4 {
			t.Errorf("unexpected click p4. expected %s but got %s\n", expectedClick.P4, click.P4)
		}

		return ch
	}

	handlerFunc := ClickHandlerFunc(expFunc)
	handlerFunc.HandleClick(context.Background(), expectedClick)
}

func TestStoreClickHandler_HandleClick(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedClick := &entity.Click{
		ID:          "expectedClickID",
		TargetURL:   "expectedTargetURL",
		TRKLink:     nil,
		UserAgent:   nil,
		CountryCode: "expectedCountryCode",
		IP:          net.ParseIP("123.1.2.3"),
		Referer:     "expectedReferer",
		P1:          "expectedP1",
		P2:          "expectedP2",
		P3:          "expectedP3",
		P4:          "expectedP4",
	}
	expectedError := errors.New("expected error")

	clkRepo := mocks.NewMockClicksRepository(ctrl)
	clkRepo.EXPECT().Save(gomock.Any(), expectedClick).Return(expectedError)

	handler := NewStoreClickHandler(clkRepo)

	outputCh := handler.HandleClick(context.Background(), expectedClick)

	clickProcessingResult := <-outputCh

	if !reflect.DeepEqual(clickProcessingResult.Click, expectedClick) {
		t.Errorf("unexpected click received from processing result. expected %+v but got %+v\n", expectedClick, clickProcessingResult.Click)
	}
	if !errors.Is(clickProcessingResult.Err, expectedError) {
		t.Errorf("unexpected error received from processing result. expected `%s` but got `%s`\n", expectedError, clickProcessingResult.Err)
	}
}
