package main

import (
	"bufio"
	"bytes"
	"fmt"
	"myiopkg"
	"os"
	"strings"
	"time"
)

func allApo(drvs []string) []string {
	var files []string
	for _, drive := range drvs {
		fmt.Printf("Поиск в папке: %s ", drive)
		file, err := myiopkg.IsApoWalkDir(drive)
		if err != nil {
			fmt.Printf("\nВНИМАНИЕ! '%s' - такой директории не существует\n", drive)
		} else if file != nil {
			fmt.Println("- найдено установленное ПО АРО3.")
		} else {
			fmt.Println("- АРО3 не обнаружено.")
		}
		files = append(files, file...)
	}
	return files
}
func argsWinPath(s string) (a2 []string) {
	a := strings.Split(s, "\\")
	if a[len(a)-1] == "" {
		a1 := make([]string, len(a)-1)
		copy(a1, a)
		a2 = []string{strings.Join(a1, "\\")}
	} else {
		a2 = []string{strings.Join(a, "\\")}
	}
	return
}
func main() {
	var (
		mapFiles map[int]string = make(map[int]string)
		drivers2 []string
	)
	drivers := myiopkg.GetDrivies()
	if len(os.Args) > 1 {
		fmt.Printf("\nЗапускаю поиск АПО3 в указанном Вами расположении: %s\n", strings.ToUpper(os.Args[1]))
		fmt.Println()
		time.Sleep(time.Second * 2)
		drivers2 = argsWinPath(strings.ToLower(os.Args[1]))
	} else {
		usrProf := os.Getenv("USERPROFILE")
		sysDrive := os.Getenv("SystemDrive")
		drivers2 = []string{sysDrive + "\\RGS", usrProf + "\\Desktop", usrProf + "\\Documents", usrProf + "\\Downloads"}
		fmt.Println("Запускаю поиск АПО3 в стандартных папках:")
		fmt.Println()
	}
	isApo := allApo(drivers2)
	fmt.Println()
	if isApo == nil {
		fmt.Println("АПО3 в стандартном размещении не найдено\nЧто делать дальше?")
		time.Sleep(time.Second * 2)
		fmt.Println("Попробовать искать АПО3 на всех доступных дисках? ВНИМАНИЕ: это может занять много времени (более 1 часа).")
		//Получить ДА или НЕТ (да, нет, д, н, yes, no, y, n - в любом регистре)
		input := myiopkg.YesNo()
		if bytes.EqualFold([]byte(input), []byte("y")) || bytes.EqualFold([]byte(input), []byte("д")) || bytes.EqualFold([]byte(input), []byte("yes")) || bytes.EqualFold([]byte(input), []byte("да")) {
			time.Sleep(time.Second * 3)
			fmt.Println("\nПоиск АПО3 по всем доступным дискам запущен...")
			isApo = allApo(drivers)
		} else {
			fmt.Println("Программа завершается...")
			time.Sleep(time.Second * 3)
			os.Exit(0)
		}
	}
	if len(isApo) == 1 {
		time.Sleep(time.Second * 1)
		fmt.Println("Найдена одна установка АРО3")
		time.Sleep(time.Second * 1)
		fmt.Println("Ищу копии базы в", "'"+isApo[0]+"\\backups\\'")
		files := myiopkg.FindAllFilesInDirByMask(isApo[0] + "\\backups\\*.bkp")
		time.Sleep(time.Second * 1)
		if files == nil {
			fmt.Println("\nФайлы для восстановления не найдены.\nВосстановить базу из сохраненных копий невозможно.")
			time.Sleep(time.Second * 1)
			fmt.Println("Копии базы в папке backups не обнаружены. Программа выключается.\n")
			time.Sleep(time.Second * 1)
			fmt.Println("Не забудьте включить в Параметрах АПО3 сохранение бэкапов базы для ее восстановления в дальнейшем.")
			time.Sleep(time.Second * 1)
			fmt.Print("Для завершения программы нажмите клавишу Enter ")
			reader := bufio.NewReader(os.Stdin)
			_, _ = reader.ReadString('\n')
			os.Exit(0)
		}
		if myiopkg.FindAndChoiceBackup(files) {
			fmt.Println("АПО3 было восстановлено на выбранную вами дату.")
			time.Sleep(time.Second * 2)
			fmt.Print("Для завершения программы нажмите клавишу Enter ")
			reader := bufio.NewReader(os.Stdin)
			_, _ = reader.ReadString('\n')
			os.Exit(0)
		}
		fmt.Println("\nПрограмма завершилась с ошибкой.")
		fmt.Println("Возможно вы не закрыли АПО. Завершите АПО и повторите восстановление еще раз.")
		fmt.Print("\nДля завершения программы нажмите клавишу Enter ")
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
		os.Exit(0)
	} else if len(isApo) > 1 {
		fmt.Println("Найдено несколько установок APO3")
		time.Sleep(time.Second * 1)
		fmt.Println("\nВыберите ту, где необходимо восстановить базу:")
		time.Sleep(time.Second * 1)
		for i, apo := range isApo {
			fmt.Printf("%d: %s\n", i+1, apo)
			mapFiles[i+1] = apo
		}
		fmt.Print("\nВыберите, в какой из установленных версий произвести восстановление \n(введите соответствующую цифру и нажмите Enter): ")
		func1 := myiopkg.ChoiceIntOpt(mapFiles)
		fmt.Printf("Выбрана версия %d: %v\n", func1, mapFiles[func1])
		time.Sleep(time.Second * 1)

		fmt.Println("Ищу копии базы в", "'"+mapFiles[func1]+"\\backups\\'")
		files := myiopkg.FindAllFilesInDirByMask(mapFiles[func1] + "\\backups\\*.bkp")
		time.Sleep(time.Second * 1)

		if myiopkg.FindAndChoiceBackup(files) {
			fmt.Println("АПО3 было восстановлено на выбранную вами дату.")
			time.Sleep(time.Second * 2)
			fmt.Print("Для завершения программы нажмите клавишу Enter ")
			reader := bufio.NewReader(os.Stdin)
			_, _ = reader.ReadString('\n')
			os.Exit(0)
		}
		fmt.Println("\nПрограмма завершилась с ошибкой.")
		fmt.Println("Возможно вы не закрыли АПО. Завершите АПО и повторите восстановление еще раз.")
		fmt.Print("\nДля завершения программы нажмите клавишу Enter ")
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
		os.Exit(0)
	} else {
		fmt.Println("\nАПО3 на компьютере не найдено.\n")
		fmt.Print("Для завершения программы нажмите клавишу Enter ")
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
		os.Exit(0)
	}
}
