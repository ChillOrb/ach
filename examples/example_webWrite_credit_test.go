package examples

import (
	"fmt"
	"github.com/moov-io/ach"
	"log"
)

func Example_webWriteCredit() {
	fh := mockFileHeader()

	bh := ach.NewBatchHeader()
	bh.ServiceClassCode = ach.CreditsOnly
	bh.CompanyName = "Name on Account"
	bh.CompanyIdentification = fh.ImmediateOrigin
	bh.StandardEntryClassCode = ach.WEB
	bh.CompanyEntryDescription = "Subscribe"
	// need EffectiveEntryDate to be fixed so it can match output
	bh.EffectiveEntryDate = "190816"
	bh.ODFIIdentification = "121042882"

	entry := ach.NewEntryDetail()
	entry.TransactionCode = ach.CheckingCredit
	entry.SetRDFI("231380104")
	entry.DFIAccountNumber = "12345678"
	entry.Amount = 10000
	entry.IndividualName = "John Doe"
	entry.SetTraceNumber(bh.ODFIIdentification, 1)
	entry.IdentificationNumber = "#789654"
	entry.DiscretionaryData = "S"
	entry.Category = ach.CategoryForward
	entry.AddendaRecordIndicator = 1

	addenda1 := ach.NewAddenda05()
	addenda1.PaymentRelatedInformation = "PAY-GATE payment\\"
	entry.AddAddenda05(addenda1)

	// build the batch
	batch := ach.NewBatchWEB(bh)
	batch.AddEntry(entry)

	if err := batch.Create(); err != nil {
		log.Fatalf("Unexpected error building batch: %s\n", err)
	}

	// build the file
	file := ach.NewFile()
	file.SetHeader(fh)
	file.AddBatch(batch)
	if err := file.Create(); err != nil {
		log.Fatalf("Unexpected error building file: %s\n", err)
	}

	fmt.Printf("%s", file.Header.String()+"\n")
	fmt.Printf("%s", file.Batches[0].GetHeader().String()+"\n")
	fmt.Printf("%s", file.Batches[0].GetEntries()[0].String()+"\n")
	fmt.Printf("%s", file.Batches[0].GetEntries()[0].Addenda05[0].String()+"\n")
	fmt.Printf("%s", file.Batches[0].GetControl().String()+"\n")
	fmt.Printf("%s", file.Control.String()+"\n")

	// Output:
	// 101 03130001202313801041908161055A094101Federal Reserve Bank   My Bank Name           12345678
	// 5220Name on Account                     231380104 WEBSubscribe       190816   1121042880000001
	// 62223138010412345678         0000010000#789654        John Doe              S 1121042880000001
	// 705PAY-GATE payment\                                                               00010000001
	// 82200000020023138010000000000000000000010000231380104                          121042880000001
	// 9000001000001000000020023138010000000000000000000010000

}
