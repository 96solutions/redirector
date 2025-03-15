package interactor_test

import (
	"context"
	"errors"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/lroman242/redirector/domain/dto"
	"github.com/lroman242/redirector/domain/entity"
	"github.com/lroman242/redirector/domain/interactor"
	"github.com/lroman242/redirector/mocks"
	"go.uber.org/mock/gomock"
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

	expFunc := func(_ context.Context, click *entity.Click) <-chan *dto.ClickProcessingResult {
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

	handlerFunc := interactor.ClickHandlerFunc(expFunc)
	handlerFunc.HandleClick(context.Background(), expectedClick)
}

func TestClickHandlerFunc_WithResult(t *testing.T) {
	testClick := &entity.Click{ID: "func-test-click"}
	expectedError := errors.New("func handler error")

	// Create a handler function that returns a specific error
	var handlerFunc interactor.ClickHandlerFunc = func(ctx context.Context, click *entity.Click) <-chan *dto.ClickProcessingResult {
		resultCh := make(chan *dto.ClickProcessingResult, 1)

		// Immediately send a result
		resultCh <- &dto.ClickProcessingResult{
			Click: click,
			Err:   expectedError,
		}
		close(resultCh)

		return resultCh
	}

	// Call the handler and verify the result
	resultCh := handlerFunc.HandleClick(context.Background(), testClick)
	result := <-resultCh

	if result.Click != testClick {
		t.Errorf("Expected Click to be %v, got %v", testClick, result.Click)
	}

	if !errors.Is(expectedError, result.Err) {
		t.Errorf("Expected Err to be %v, got %v", expectedError, result.Err)
	}
}

func TestNewStoreClickHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test that NewStoreClickHandler returns a correctly initialized handler
	mockRepo := mocks.NewMockClicksRepository(ctrl)
	handler := interactor.NewStoreClickHandler(mockRepo)

	// Since storeClickHandler is unexported, just check we get a non-nil interface
	if handler == nil {
		t.Error("NewStoreClickHandler returned nil")
	}

	// We'll need to expect a Save call if we're going to test the handler
	testClick := &entity.Click{ID: "test-constructor"}
	mockRepo.EXPECT().
		Save(gomock.Any(), testClick).
		Return(nil)

	// Verify we can call methods on the handler
	outputCh := handler.HandleClick(context.Background(), testClick)
	if outputCh == nil {
		t.Error("Handler returned nil channel")
	}

	// Make sure we can read from the channel
	result := <-outputCh
	if result.Click != testClick {
		t.Error("Unexpected click in result")
	}
}

func TestStoreClickHandler_HandleClick_WithError(t *testing.T) {
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

	handler := interactor.NewStoreClickHandler(clkRepo)

	outputCh := handler.HandleClick(context.Background(), expectedClick)

	clickProcessingResult := <-outputCh

	if !reflect.DeepEqual(clickProcessingResult.Click, expectedClick) {
		t.Errorf("unexpected click received from processing result. expected %+v but got %+v\n",
			expectedClick, clickProcessingResult.Click)
	}
	if !errors.Is(clickProcessingResult.Err, expectedError) {
		t.Errorf("unexpected error received from processing result. expected `%s` but got `%s`\n",
			expectedError, clickProcessingResult.Err)
	}
}

func TestStoreClickHandler_HandleClick_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedClick := &entity.Click{
		ID:          "successClickID",
		TargetURL:   "successTargetURL",
		TRKLink:     nil,
		UserAgent:   nil,
		CountryCode: "US",
		IP:          net.ParseIP("192.168.1.1"),
		Referer:     "https://example.com",
		P1:          "p1val",
		P2:          "p2val",
		P3:          "p3val",
		P4:          "p4val",
	}

	clkRepo := mocks.NewMockClicksRepository(ctrl)
	// Repository returns nil error (success case)
	clkRepo.EXPECT().Save(gomock.Any(), expectedClick).Return(nil)

	handler := interactor.NewStoreClickHandler(clkRepo)

	outputCh := handler.HandleClick(context.Background(), expectedClick)

	clickProcessingResult := <-outputCh

	if !reflect.DeepEqual(clickProcessingResult.Click, expectedClick) {
		t.Errorf("unexpected click received from processing result. expected %+v but got %+v\n",
			expectedClick, clickProcessingResult.Click)
	}
	if clickProcessingResult.Err != nil {
		t.Errorf("unexpected error received from processing result. expected nil but got `%s`\n",
			clickProcessingResult.Err)
	}
}

func TestStoreClickHandler_HandleClick_ContextCancellation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedClick := &entity.Click{
		ID:          "cancelledClickID",
		TargetURL:   "cancelledTargetURL",
		TRKLink:     nil,
		UserAgent:   nil,
		CountryCode: "UK",
		IP:          net.ParseIP("10.0.0.1"),
		Referer:     "https://example.org",
		P1:          "p1test",
		P2:          "p2test",
		P3:          "p3test",
		P4:          "p4test",
	}

	// Create a context that will be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	clkRepo := mocks.NewMockClicksRepository(ctrl)
	// No need to expect Save() call as it should be cancelled before reaching that point

	handler := interactor.NewStoreClickHandler(clkRepo)

	// Get the channel first
	outputCh := handler.HandleClick(ctx, expectedClick)

	// Cancel the context immediately
	cancel()

	// Get the result which should contain a context cancellation error
	clickProcessingResult := <-outputCh

	if !reflect.DeepEqual(clickProcessingResult.Click, expectedClick) {
		t.Errorf("unexpected click received from processing result. expected %+v but got %+v\n",
			expectedClick, clickProcessingResult.Click)
	}
	if clickProcessingResult.Err == nil {
		t.Error("expected context cancellation error, but got nil")
	} else if !errors.Is(clickProcessingResult.Err, context.Canceled) {
		t.Errorf("expected context.Canceled error, but got %v", clickProcessingResult.Err)
	}
}

func TestStoreClickHandler_HandleClick_ContextTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedClick := &entity.Click{
		ID:          "timeoutClickID",
		TargetURL:   "timeoutTargetURL",
		CountryCode: "DE",
		IP:          net.ParseIP("172.16.0.1"),
	}

	// Create a context that will time out
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	clkRepo := mocks.NewMockClicksRepository(ctrl)
	// No need to expect Save() call as it should timeout before reaching that point

	handler := interactor.NewStoreClickHandler(clkRepo)

	// Sleep to ensure timeout occurs
	time.Sleep(5 * time.Millisecond)

	// Get the channel first
	outputCh := handler.HandleClick(ctx, expectedClick)

	// Get the result which should contain a context deadline exceeded error
	clickProcessingResult := <-outputCh

	if !reflect.DeepEqual(clickProcessingResult.Click, expectedClick) {
		t.Errorf("unexpected click received from processing result. expected %+v but got %+v\n",
			expectedClick, clickProcessingResult.Click)
	}
	if clickProcessingResult.Err == nil {
		t.Error("expected context deadline exceeded error, but got nil")
	} else if clickProcessingResult.Err != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded error, but got %v", clickProcessingResult.Err)
	}
}
