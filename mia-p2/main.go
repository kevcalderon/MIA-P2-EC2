package main

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// Main function
func main() {
	ListaDiscos := list.New()
	LlenarListaDisco(ListaDiscos)
	//LeerArchivo de Entrada
	var comando string = ""
	scanner := bufio.NewScanner(os.Stdin)
	for comando != "exit" {
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Println("                            MIA PROYECTO2 EC2                                   ")
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Print(">>")
		scanner.Scan()
		comando = scanner.Text()
		if comando != "" {
			LeerTexto(comando, ListaDiscos)
		}

	}
}

func LlenarListaDisco(ListaDiscos *list.List) {
	IdDisco := [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i",
		"j", "k", "l", "m", "n", "o", "p", "q",
		"r", "s", "t", "u", "v", "w", "x", "y", "z"}
	for i := 0; i < 26; i++ {
		disco := DISCO{}
		copy(disco.Estado[:], "0")
		copy(disco.Id[:], IdDisco[i])
		for j := 0; j < len(disco.Particiones); j++ {
			mount := MOUNT{}
			mount.NombreParticion = ""
			mount.Id = strconv.Itoa(j + 1)
			copy(mount.Estado[:], "0")
			disco.Particiones[j] = mount
		}
		ListaDiscos.PushBack(disco)
	}
}

var global string = ""

//Funcion para leer y reconocer los comandos lleno la lista de comandos
func LeerTexto(dat string, ListaDiscos *list.List) {
	//Leendo la cadena de entrada
	ListaComandos := list.New()
	lineaComando := strings.Split(dat, "\n")
	var c Comando
	for i := 0; i < len(lineaComando); i++ {
		EsComentario := lineaComando[i][0:1]
		if EsComentario != "#" {
			comando := lineaComando[i]
			if strings.Contains(lineaComando[i], "\\*") {
				comando = strings.Replace(lineaComando[i], "\\*", " ", 1) + lineaComando[i+1]
				i = i + 1
			}
			propiedades := strings.Split(string(comando), " ")
			nombreComando := propiedades[0]
			c = Comando{Name: strings.ToLower(nombreComando)}
			propiedadesTemp := make([]Propiedad, len(propiedades)-1)
			for i := 1; i < len(propiedades); i++ {
				if propiedades[i] == "" {
					continue
				} else if propiedades[i] == "-p" {
					propiedadesTemp[i-1] = Propiedad{Name: "-p",
						Val: "-p"}
				} else {
					if strings.Contains(propiedades[i], "=") {
						valor_propiedad_Comando := strings.Split(propiedades[i], "=")
						propiedadesTemp[i-1] = Propiedad{Name: valor_propiedad_Comando[0],
							Val: valor_propiedad_Comando[1]}
					} else {
						propiedadesTemp[i-1] = Propiedad{Name: "-sigue",
							Val: propiedades[i]}
					}
				}
			}
			c.Propiedades = propiedadesTemp
			ListaComandos.PushBack(c)
		}
	}
	RecorrerListaComando(ListaComandos, ListaDiscos)
}

func RecorrerListaComando(ListaComandos *list.List, ListaDiscos *list.List) {
	var ParamValidos bool = true
	var cont = 1
	for element := ListaComandos.Front(); element != nil; element = element.Next() {
		var comandoTemp Comando
		comandoTemp = element.Value.(Comando)
		switch strings.ToLower(comandoTemp.Name) {
		case "mkdisk":
			ParamValidos = EjecutarComandoMKDISK(comandoTemp.Name, comandoTemp.Propiedades, cont)
			cont++
			if ParamValidos == false {
				fmt.Println("Parametros Invalidos")
			}
		case "rmdisk":
			ParamValidos = EjecutarComandoRMDISK(comandoTemp.Name, comandoTemp.Propiedades)
			if ParamValidos == false {
				fmt.Println("Parametros Invalidos")
			}
		case "fdisk":
			ParamValidos = EjecutarComandoFDISK(comandoTemp.Name, comandoTemp.Propiedades)
			if ParamValidos == false {
				fmt.Println("Parametros Invalidos")
			}
		case "mount":
			if len(comandoTemp.Propiedades) != 0 {
				ParamValidos = EjecutarComandoMount(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
				if ParamValidos == false {
					fmt.Println("Parametros Invalidos")
				} else {
					EjecutarReporteMount(ListaDiscos)
				}
			}
		case "exit":
			fmt.Println("")
		case "pause":
			fmt.Println("Presione la tecla enter para Continuar")
			fmt.Scanln()
		case "exec":
			ParamValidos = EjecutarComandoExec(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("Error, parametros no validos")
			}
		case "mkdir":
			//ParamValidos = EjecutarComandoMKDIR(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("Error, parametros no validos")
			}
		case "mkfile":
			//ParamValidos = EjecutarComandoMKFILE(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("Error, parametros no validos")
			}
		case "mkfs":
			ParamValidos = EjecutarComandoMKFS(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("Error, parametros no validos")
			}
		case "rep":
			ParamValidos = EjecutarComandoReporte(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("Error, parametros no validos")
			}
		case "login":
			ParamValidos, global = EjecutarComandoLogin(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)

			if ParamValidos == false {
				fmt.Println("Error, parametros no validos")
			}
		case "logout":
			/*if global == "" {
				fmt.Println("Inicia sesión primero para realizar la acción.")
			} else {
				global = ""
				fmt.Println("Sesión finalizada")
			}*/
		default:
			fmt.Println("Error, comando invalido")
		}
	}
}

/*----------------------RMDISK------------------------------------------------*/
func EjecutarComandoRMDISK(nombreComando string, propiedadesTemp []Propiedad) (ParamValidos bool) {
	ParamValidos = true
	var opciones string = ""
	scanner := bufio.NewScanner(os.Stdin)
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-path":
				fmt.Println("¿Estas seguro eliminar el disco? S/N")
				scanner.Scan()
				opciones = scanner.Text()
				if opciones == "S" {
					executeComand("rm " + propiedadTemp.Val)
					fmt.Println("Se ha eliminado el disco correctamente")
				}
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

/*------------------------------MKDISK-------------------------------------------*/
func EjecutarComandoMKDISK(nombreComando string, propiedadesTemp []Propiedad, cont int) (ParamValidos bool) {
	dt := time.Now()
	mbr1 := MBR{}
	copy(mbr1.MbrFechaCreacion[:], dt.String())
	mbr1.NoIdentificador = int64(rand.Intn(100) + cont)
	fmt.Println("----------------- Ejecutando MKDISK -----------------")
	comandos := "dd if=/dev/zero "
	ParamValidos = true
	pathCompleta := ""
	var propiedades [4]string
	if len(propiedadesTemp) >= 2 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			//Vector temporal de propiedades
			switch strings.ToLower(nombrePropiedad) {
			case "-size":
				propiedades[0] = propiedadTemp.Val
			case "-unit":
				propiedades[2] = strings.ToLower(propiedadTemp.Val)
			case "-path":
				propiedades[3] = propiedadTemp.Val
				arr_path := strings.Split(propiedades[3], "/")

				for i := 0; i < len(arr_path)-1; i++ {
					pathCompleta += arr_path[i] + "/"
				}
				executeComand("mkdir " + pathCompleta)
				comandos += "of=" + propiedades[3]
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		EsComilla := propiedades[3][0:1]
		if EsComilla == "\"" {
			propiedades[3] = propiedades[3][1 : len(propiedades[3])-1]
		}
		tamanioTotal, _ := strconv.ParseInt(propiedades[0], 10, 64)
		if propiedades[2] == "k" {
			comandos += " bs=" + strconv.Itoa((int(tamanioTotal))*1000) + " count=1"
			mbr1.MbrTamanio = ((tamanioTotal) - 1) * 1000
		} else {
			comandos += " bs=" + strconv.Itoa(int(tamanioTotal)) + "MB" + " count=1"
			mbr1.MbrTamanio = tamanioTotal * 1000000
		}
		//Inicializando Particiones
		for i := 0; i < 4; i++ {
			copy(mbr1.Particiones[i].Status_particion[:], "0")
			copy(mbr1.Particiones[i].TipoParticion[:], "")
			copy(mbr1.Particiones[i].TipoAjuste[:], "")
			mbr1.Particiones[i].Inicio_particion = 0
			mbr1.Particiones[i].TamanioTotal = 0
			copy(mbr1.Particiones[i].NombreParticion[:], "")
		}

		executeComand(comandos)

		f, err := os.OpenFile(propiedades[3], os.O_WRONLY, 0755)
		if err != nil {
			log.Fatalln(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalln(err)
			}
		}()
		f.Seek(0, 0)
		err = binary.Write(f, binary.BigEndian, mbr1)
		if err != nil {
			log.Fatalln(err, propiedades[3])
		}
		fmt.Println("Se ha creado el disco en la carpeta: ", pathCompleta)
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

/*-----------------------EXEC----------------------------------------------*/
func EjecutarComandoExec(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	fmt.Println("----------------- Ejecutando EXEC -----------------")
	ParamValidos = true
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-path":
				fmt.Println(propiedadTemp.Val)
				dat, err := ioutil.ReadFile(propiedadTemp.Val)
				CheckError(err)
				LeerTexto(string(dat), ListaDiscos)
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

func executeComand(comandos string) {
	args := strings.Split(comandos, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.CombinedOutput()
}
func BytesToString(data [1]byte) string {
	return string(data[:])
}
func CheckError(e error) {
	if e != nil {
		fmt.Println("Error - ----------")
		fmt.Println(e)
	}
}

/*----------------------------FDISK---------------------------------------*/
func EjecutarComandoFDISK(nombreComando string, propiedadesTemp []Propiedad) (ParamValidos bool) {
	fmt.Println("----------------- Ejecutando FDISK -----------------")
	ParamValidos = true
	mbr := MBR{}
	particion := Particion{}
	var startPart int64 = int64(unsafe.Sizeof(mbr))
	var propiedades [8]string
	if len(propiedadesTemp) >= 2 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-size":
				propiedades[0] = propiedadTemp.Val
			case "-fit":
				propiedades[1] = propiedadTemp.Val
			case "-unit":
				propiedades[2] = propiedadTemp.Val
			case "-path":
				propiedades[3] = propiedadTemp.Val
			case "-type":
				propiedades[4] = propiedadTemp.Val
			case "-delete":
				propiedades[5] = propiedadTemp.Val
			case "-name":
				propiedades[6] = propiedadTemp.Val
				fmt.Println(propiedades[6])
			case "-add":
				propiedades[7] = propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		EsComilla := propiedades[3][0:1]
		if EsComilla == "\"" {
			propiedades[3] = propiedades[3][1 : len(propiedades[3])-1]
		}
		//Tamanio Particion
		var TamanioTotalParticion int64 = 0
		if strings.ToLower(propiedades[2]) == "b" {
			TamanioParticion, _ := strconv.ParseInt(propiedades[0], 10, 64)
			TamanioTotalParticion = TamanioParticion
		} else if strings.ToLower(propiedades[2]) == "k" {
			TamanioParticion, _ := strconv.ParseInt(propiedades[0], 10, 64)
			TamanioTotalParticion = TamanioParticion * 1000
		} else if strings.ToLower(propiedades[2]) == "m" {
			TamanioParticion, _ := strconv.ParseInt(propiedades[0], 10, 64)
			TamanioTotalParticion = TamanioParticion * 1000000
		} else {
			TamanioParticion, _ := strconv.ParseInt(propiedades[0], 10, 64)
			TamanioTotalParticion = TamanioParticion * 1000
		}
		if propiedades[5] != "" {
			EliminarParticion(propiedades[3], propiedades[6], propiedades[5])
			return
		}
		//Obtener el MBR
		switch strings.ToLower(propiedades[4]) {
		case "p":
			var Particiones [4]Particion
			f, err := os.OpenFile(propiedades[3], os.O_RDWR, 0755)
			if err != nil {
				fmt.Println("No existe la ruta" + propiedades[3])
				return false
			}
			defer f.Close()
			f.Seek(0, 0)
			err = binary.Read(f, binary.BigEndian, &mbr)
			Particiones = mbr.Particiones
			if err != nil {
				fmt.Println("No existe el archivo en la ruta")
			}

			if HayEspacio(TamanioTotalParticion, mbr.MbrTamanio) {
				return false
			} //Verificar si ya hay particiones
			if BytesToString(Particiones[0].Status_particion) == "1" {
				for i := 0; i < 4; i++ {
					//Posicion en bytes del partstar de la n particion
					startPart += Particiones[i].TamanioTotal
					if BytesToString(Particiones[i].Status_particion) == "0" {
						//fmt.Println(startPart)
						break
					}
				}
			}
			if HayEspacio(startPart+TamanioTotalParticion, mbr.MbrTamanio) {
				return false
			}
			//dando valores a la particion
			copy(particion.Status_particion[:], "1")
			copy(particion.TipoParticion[:], propiedades[4])
			copy(particion.TipoAjuste[:], propiedades[1])
			particion.Inicio_particion = startPart
			particion.TamanioTotal = TamanioTotalParticion
			copy(particion.NombreParticion[:], propiedades[6])
			//Particion creada
			for i := 0; i < 4; i++ {
				if BytesToString(Particiones[i].Status_particion) == "0" {
					Particiones[i] = particion
					break
				}
			}
			f.Seek(0, 0)
			mbr.Particiones = Particiones
			err = binary.Write(f, binary.BigEndian, mbr)
			ReadFile(propiedades[3])
		case "l":
			fmt.Println("Particion Logica")
			if !HayExtendida(propiedades[3]) {
				fmt.Println("No existe una particion Extendida")
				return false
			}
			ebr := EBR{}
			copy(ebr.Status_particion[:], "1")
			copy(ebr.TipoAjuste[:], propiedades[1])
			ebr.Inicio_particion = startPart
			ebr.Particion_Siguiente = 0
			ebr.TamanioTotal = TamanioTotalParticion
			copy(ebr.NombreParticion[:], propiedades[6])
			//Obteniendo el byte donde empezara la particion Logica
			InicioParticionLogica(propiedades[3], ebr)
		case "e":
			//Particiones Extendidas
			var Particiones [4]Particion
			f, err := os.OpenFile(propiedades[3], os.O_RDWR, 0755)
			if err != nil {
				fmt.Println("No existe la ruta" + propiedades[3])
				return false
			}
			defer f.Close()
			f.Seek(0, 0)
			err = binary.Read(f, binary.BigEndian, &mbr)
			Particiones = mbr.Particiones
			if err != nil {
				fmt.Println("No existe el archivo en la ruta")
			}
			//El mbr ya se a leido,2.Verificar si existe espacion disponible o que no lo rebase
			if HayEspacio(TamanioTotalParticion, mbr.MbrTamanio) {
				return false
			} //Verificar si ya hay particiones
			if BytesToString(Particiones[0].Status_particion) == "1" {
				for i := 0; i < 4; i++ {
					//Posicion en bytes del partstar de la n particion
					startPart += Particiones[i].TamanioTotal
					if BytesToString(Particiones[i].Status_particion) == "0" {

						break
					}
				}
			}
			if HayEspacio(startPart+TamanioTotalParticion, mbr.MbrTamanio) {
				return false
			}
			//dando valores a la particion
			copy(particion.Status_particion[:], "1")
			copy(particion.TipoParticion[:], propiedades[4])
			copy(particion.TipoAjuste[:], propiedades[1])
			particion.Inicio_particion = startPart
			particion.TamanioTotal = TamanioTotalParticion
			copy(particion.NombreParticion[:], propiedades[6])
			//Particion creada
			for i := 0; i < 4; i++ {
				if BytesToString(Particiones[i].Status_particion) == "0" {
					Particiones[i] = particion
					break
				}
			}
			f.Seek(0, 0)
			mbr.Particiones = Particiones
			err = binary.Write(f, binary.BigEndian, mbr)
			ReadFile(propiedades[3])
			ebr := EBR{}
			copy(ebr.Status_particion[:], "1")
			copy(ebr.TipoAjuste[:], propiedades[1])
			ebr.Inicio_particion = startPart
			ebr.Particion_Siguiente = -1
			ebr.TamanioTotal = TamanioTotalParticion
			copy(ebr.NombreParticion[:], propiedades[6])
			f.Seek(ebr.Inicio_particion, 0)
			err = binary.Write(f, binary.BigEndian, ebr)

			fmt.Println("Extendida", "Leendo EBR")
		default:
			fmt.Println("Ocurrio un error")
		}
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}
func EscribirParticionLogica(path string, ebr EBR, inicioParticionLogica int64, inicioParticionExtendida int64) bool {
	ebr.Inicio_particion = inicioParticionLogica
	ebr.Particion_Siguiente = inicioParticionLogica + ebr.TamanioTotal
	return true
}

func EliminarParticion(path string, name string, typeDelete string) bool {
	var name2 [15]byte
	Encontrada := false
	copy(name2[:], name)
	f, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	//Posiciono al inicio el Puntero
	f.Seek(0, 0)
	//Leo el mbr
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" && BytesNombreParticion(Particiones[i].NombreParticion) == BytesNombreParticion(name2) {
			fmt.Println("Es una Extendida")
			Encontrada = true
		} else if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "p" && BytesNombreParticion(Particiones[i].NombreParticion) == BytesNombreParticion(name2) {
			var partTemp = Particion{}
			copy(partTemp.Status_particion[:], "0")
			copy(partTemp.TipoParticion[:], "")
			copy(partTemp.TipoAjuste[:], "")
			partTemp.Inicio_particion = 0
			partTemp.TamanioTotal = 0
			copy(partTemp.NombreParticion[:], "")
			Particiones[i] = partTemp
			mbr.Particiones = Particiones
			f.Seek(0, 0)
			err = binary.Write(f, binary.BigEndian, &mbr)
			fmt.Println("Particion Primaria Eliminada")
			ReadFile(path)
			Encontrada = true
		}
	}
	if Encontrada == false {
		for i := 0; i < 4; i++ {
			if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
				var InicioExtendida int64 = Particiones[i].Inicio_particion
				f.Seek(InicioExtendida, 0)
				ebrAnterior := EBR{}
				ebr := EBR{}
				ebrAnterior = ebr
				err = binary.Read(f, binary.BigEndian, &ebr)
				if ebr.Particion_Siguiente == -1 {
					fmt.Println("No Hay particiones Logicas")
				} else {
					f.Seek(InicioExtendida, 0)
					err = binary.Read(f, binary.BigEndian, &ebr)
					for {
						if BytesNombreParticion(ebr.NombreParticion) == BytesNombreParticion(name2) {
							fmt.Println("Particion Logica Encontrada")
							if strings.ToLower(typeDelete) == "fast" {
								ebrAnterior.Particion_Siguiente = ebr.Particion_Siguiente
								f.Seek(ebrAnterior.Inicio_particion, 0)
								err = binary.Write(f, binary.BigEndian, ebrAnterior)

							} else if strings.ToLower(typeDelete) == "full" {
								ebrAnterior.Particion_Siguiente = ebr.Particion_Siguiente
								f.Seek(ebrAnterior.Inicio_particion, 0)
								err = binary.Write(f, binary.BigEndian, ebrAnterior)
							}
							Encontrada = true
						}
						if ebr.Particion_Siguiente == -1 {
							break
						} else {
							f.Seek(ebr.Particion_Siguiente, 0)
							ebrAnterior = ebr
							err = binary.Read(f, binary.BigEndian, &ebr)
						}
					}
				}
			}
		}
	}
	if Encontrada == false {
		fmt.Println("Error No se encontro la particon")
	}
	return false
}
func InicioParticionLogica(path string, ebr2 EBR) bool {
	f, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
			var InicioExtendida int64 = Particiones[i].Inicio_particion
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Particion_Siguiente == -1 {
				ebr.Particion_Siguiente = ebr.Inicio_particion + int64(unsafe.Sizeof(ebr)) + ebr2.TamanioTotal
				f.Seek(InicioExtendida, 0)
				err = binary.Write(f, binary.BigEndian, ebr)
				ebr2.Inicio_particion = ebr.Particion_Siguiente
				ebr2.Particion_Siguiente = -1
				f.Seek(ebr2.Inicio_particion, 0)
				err = binary.Write(f, binary.BigEndian, ebr2)

				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				fmt.Println(ebr.Inicio_particion)
				fmt.Println(ebr.Particion_Siguiente)
				return false
			} else {

				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Particion_Siguiente == -1 {
						//fmt.Println("Es la ultima logica")
						ebr.Particion_Siguiente = ebr.Inicio_particion + int64(unsafe.Sizeof(ebr)) + ebr2.TamanioTotal
						f.Seek(ebr.Inicio_particion, 0)
						err = binary.Write(f, binary.BigEndian, ebr)
						ebr2.Inicio_particion = ebr.Particion_Siguiente
						ebr2.Particion_Siguiente = -1
						f.Seek(ebr2.Inicio_particion, 0)
						err = binary.Write(f, binary.BigEndian, ebr2)
						fmt.Printf("NombreLogica: %s\n", ebr2.NombreParticion)
						break
					} else {
						f.Seek(ebr.Particion_Siguiente, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
						fmt.Printf("NombreLogica: %s\n", ebr.NombreParticion)
					}

				}
				return false
			}
		}
	}
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}

	return false
}
func HayExtendida(path string) bool {
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
			return true
		}
	}
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}

	return false
}
func ReadFileEBR(path string) (funciona bool) {
	//fmt.Println("****************Leendo EL EBR")
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
			var InicioExtendida int64 = Particiones[i].Inicio_particion
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Particion_Siguiente == -1 {
				fmt.Println("No Hay particiones Logicas")
			} else {
				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Particion_Siguiente == -1 {
						break
					} else {
						f.Seek(ebr.Particion_Siguiente, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
					}
					fmt.Printf("NombreLogica: %s\n", ebr.NombreParticion)

				}
			}
		}
	}
	return true
}
func ReadFile(path string) (funciona bool) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	return true
}
func HayEspacio(TamanioTotalParticion int64, tamanioDisco int64) bool {
	if ((TamanioTotalParticion) > tamanioDisco) || (TamanioTotalParticion < 0) {
		fmt.Println("Error Fdisk, El Tamanio de la particion es mayor a el tamanio del disco o el tamanio es incorrecto")
		return true
	}
	return false
}

/*------------------------------MOUNT----------------------------------------*/
func EjecutarComandoMount(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	var propiedades [2]string
	var nombre [15]byte
	ParamValidos = true
	if len(propiedadesTemp) >= 2 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-name":
				propiedades[0] = propiedadTemp.Val
				copy(nombre[:], propiedades[0])
			case "-path":
				propiedades[1] = propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		//Empezar a montar las Particiones
		EjecutarComando(propiedades[1], nombre, ListaDiscos)
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}
func EjecutarReporteMount(ListaDiscos *list.List) {
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if disco.NombreDisco != "" {
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if mountTemp.NombreParticion != "" {
					fmt.Println("id=", mountTemp.Id, " ruta=", disco.Path, "nombre_particion=", mountTemp.NombreParticion)
				}
			}
		}
	}
}
func IdValido(id string, ListaDiscos *list.List) bool {
	esta := false
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if disco.NombreDisco != "" {
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if mountTemp.NombreParticion != "" {
					if mountTemp.Id == id {
						return true
					}
				}
			}
		}
	}
	return esta
}

func EjecutarComando(path string, NombreParticion [15]byte, ListaDiscos *list.List) bool {
	var encontrada = false
	lineaComando := strings.Split(path, "/")
	nombreDisco := lineaComando[len(lineaComando)-1]
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	for i := 0; i < 4; i++ {
		if string(Particiones[i].NombreParticion[:]) == string(NombreParticion[:]) {
			encontrada = true
			if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
				fmt.Println("Error no se puede Montar una particion Extendida")
			} else {
				ParticionMontar(ListaDiscos, string(NombreParticion[:]), string(nombreDisco), path)
			}
		}
		if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
			ebr := EBR{}
			f.Seek(Particiones[i].Inicio_particion, 0)
			err = binary.Read(f, binary.BigEndian, &ebr)
			for {
				if ebr.Particion_Siguiente == -1 {
					break
				} else {
					f.Seek(ebr.Particion_Siguiente, 0)
					err = binary.Read(f, binary.BigEndian, &ebr)
				}
				var nombre string = string(ebr.NombreParticion[:])
				var nombre2 string = string(NombreParticion[:])
				if nombre == nombre2 {
					encontrada = true
					//Montar Particion
					ParticionMontar(ListaDiscos, string(NombreParticion[:]), string(nombreDisco), path)
				}
			}
		}
	}
	if encontrada == false {
		fmt.Println("Error no se encontro la particion")
	}
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	return true
}

func ParticionMontar(ListaDiscos *list.List, nombreParticion string, nombreDisco string, path string) {

	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if BytesToString(disco.Estado) == "0" && !ExisteDisco(ListaDiscos, nombreDisco) {
			disco.NombreDisco = nombreDisco
			disco.Path = path
			copy(disco.Estado[:], "1")
			//#id->vda1
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if BytesToString(mountTemp.Estado) == "0" {
					mountTemp.Id = "14" + BytesToString(disco.Id) + mountTemp.Id
					mountTemp.NombreParticion = nombreParticion
					copy(mountTemp.Estado[:], "1")
					copy(mountTemp.EstadoMKS[:], "0")
					disco.Particiones[i] = mountTemp
					break
				} else if BytesToString(mountTemp.Estado) == "1" && mountTemp.NombreParticion == nombreParticion {
					break
				}
			}
			element.Value = disco
			break
		} else if BytesToString(disco.Estado) == "1" && ExisteDisco(ListaDiscos, nombreDisco) && nombreDisco == disco.NombreDisco {

			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if BytesToString(mountTemp.Estado) == "0" {
					mountTemp.Id = "14" + BytesToString(disco.Id) + mountTemp.Id
					mountTemp.NombreParticion = nombreParticion
					copy(mountTemp.Estado[:], "1")
					copy(mountTemp.EstadoMKS[:], "0")
					disco.Particiones[i] = mountTemp
					break
				} else if BytesToString(mountTemp.Estado) == "1" && mountTemp.NombreParticion == nombreParticion {
					//fmt.Println("La Particion ya esta montada")
					break
				}
			}
			element.Value = disco
			break
		}
	}
}
func ExisteDisco(ListaDiscos *list.List, nombreDisco string) bool {
	Existe := false
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if disco.NombreDisco == nombreDisco {
			return true
		} else {
			Existe = false
		}
	}
	return Existe
}

/*--------------------------REPORTES----------------------------------------*/
func EjecutarComandoReporte(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	ParamValidos = true
	var propiedades [4]string
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-id":
				propiedades[0] = propiedadTemp.Val
			case "-path":
				propiedades[1] = propiedadTemp.Val
			case "-name":
				propiedades[2] = propiedadTemp.Val
			case "-ruta":
				propiedades[3] = propiedadTemp.Val
			case "-sigue":
				propiedades[1] += propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando", nombrePropiedad)
			}
		}
		EsComilla := propiedades[1][0:1]
		if EsComilla == "\"" {
			if propiedades[3] != "" {
				propiedades[3] = propiedades[3][1 : len(propiedades[3])-1]
			}
			propiedades[1] = propiedades[1][1 : len(propiedades[1])-1]
		}
		carpetas_Graficar := strings.Split(propiedades[1], "/")
		var comando = ""
		for i := 1; i < len(carpetas_Graficar)-1; i++ {
			comando += carpetas_Graficar[i] + "/"
		}
		fmt.Println(comando)
		executeComand("mkdir " + comando[0:len(comando)-1])
		switch strings.ToLower(propiedades[2]) {
		case "mbr":
			GraficarDisco(propiedades[0], ListaDiscos, propiedades[1])
		case "tree":
			//GraficarTreeFull(propiedades[0], propiedades[1], propiedades[3], ListaDiscos)
		default:
			fmt.Println("Reporte incorrecto.")

		}
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

func GraficarDisco(idParticion string, ListaDiscos *list.List, path string) bool {
	var NombreParticion [15]byte
	var buffer bytes.Buffer
	buffer.WriteString("digraph G{\ntbl [\nshape=box\nlabel=<\n<table border='0' cellborder='0' width='100' height=\"30\">\n<tr>")
	pathDisco, _, _ := RecorrerListaDisco(idParticion, ListaDiscos)
	fmt.Println(pathDisco)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	PorcentajeUtilizao := 0.0
	var EspacioUtilizado int64 = 0
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	TamanioDisco := mbr.MbrTamanio
	Particiones := mbr.Particiones
	buffer.WriteString("<td height='30' width='75'> MBR </td>")
	for i := 0; i < 4; i++ {
		if convertName(Particiones[i].NombreParticion[:]) != convertName(NombreParticion[:]) && strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "p" {
			PorcentajeUtilizao = (float64(Particiones[i].TamanioTotal) / float64(TamanioDisco)) * 100
			buffer.WriteString("<td height='30' width='75.0'>PRIMARIA <br/>" + convertName(Particiones[i].NombreParticion[:]) + " <br/> Ocupado: " + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>")
			EspacioUtilizado += Particiones[i].TamanioTotal
		} else if convertName(Particiones[i].Status_particion[:]) == "0" {
			buffer.WriteString("<td height='30' width='75.0'>Libre</td>")
		}
		if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
			EspacioUtilizado += Particiones[i].TamanioTotal
			PorcentajeUtilizao = (float64(Particiones[i].TamanioTotal) / float64(TamanioDisco)) * 100
			buffer.WriteString("<td  height='30' width='15.0'>\n")
			buffer.WriteString("<table border='5'  height='30' WIDTH='15.0' cellborder='1'>\n")
			buffer.WriteString(" <tr>  <td height='60' colspan='100%'>EXTENDIDA <br/>" + convertName(Particiones[i].NombreParticion[:]) + " <br/> Ocupado:" + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>  </tr>\n<tr>")
			var InicioExtendida int64 = Particiones[i].Inicio_particion
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Particion_Siguiente == -1 {
				fmt.Println("No Hay particiones Logicas")
			} else {
				var EspacioUtilizado int64 = 0
				cont := 0
				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Particion_Siguiente == -1 {
						break
					} else {
						f.Seek(ebr.Particion_Siguiente, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
						EspacioUtilizado += ebr.TamanioTotal
						PorcentajeUtilizao = (float64(ebr.TamanioTotal) / float64(Particiones[i].TamanioTotal)) * 100
						buffer.WriteString("<td height='30'>EBR</td><td height='30'> Logica:  " + convertName(ebr.NombreParticion[:]) + " " + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>")
						cont++
					}
				}
				if (Particiones[i].TamanioTotal - EspacioUtilizado) > 0 {
					PorcentajeUtilizao = (float64(TamanioDisco-EspacioUtilizado) / float64(TamanioDisco)) * 100
					buffer.WriteString("<td height='30' width='100%'>Libre: " + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>")
				}
			}
			buffer.WriteString("</tr>\n")
			buffer.WriteString("</table>\n</td>")
		}
	}
	if (TamanioDisco - EspacioUtilizado) > 0 {
		PorcentajeUtilizao = (float64(TamanioDisco-EspacioUtilizado) / float64(TamanioDisco)) * 100
		buffer.WriteString("<td height='30' width='75.0'>Libre: " + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>")
	}
	buffer.WriteString("     </tr>\n</table>\n>];\n}")
	var datos string
	datos = string(buffer.String())
	CreateArchivo(path, datos)
	return false
}

func CreateArchivo(path string, data string) {
	propiedades := strings.Split(path, "/")
	nombreArchivo := propiedades[len(propiedades)-1]
	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(data)

	if err2 != nil {
		log.Fatal(err2)
	}
	executeComand("dot -Tpdf " + path + " -o " + nombreArchivo[0:len(nombreArchivo)-4] + ".pdf")
	//executeComand("xdg-open " + nombreArchivo[0:len(nombreArchivo)-4] + ".pdf")
	//executeComand("xdg-open " + path)

}

func convertName(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

/*-----------------MKFS------------------------------------ */
func EjecutarComandoMKFS(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	ParamValidos = true
	var propiedades [4]string
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-id":
				propiedades[0] = propiedadTemp.Val
			case "-type":
				propiedades[1] = propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		ExecuteMKFS(propiedades[0], ListaDiscos)
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}
func ExecuteMKFS(id string, ListaDiscos *list.List) bool {
	dt := time.Now()
	idValido := IdValido(id, ListaDiscos)
	if idValido == false {
		fmt.Println("El id no existe")
		return false
	}
	Id := strings.ReplaceAll(id, "14", "")
	NoParticion := Id[1:]
	IdDisco := Id[:1]
	pathDisco := ""
	nombreParticion := ""
	nombreDisco := ""
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if BytesToString(disco.Id) == IdDisco {
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if mountTemp.Id == id {
					copy(mountTemp.EstadoMKS[:], "1")
					nombreParticion = mountTemp.NombreParticion
					pathDisco = disco.Path
					nombreDisco = disco.NombreDisco
					break
				}
			}

		}
		element.Value = disco
	}
	mbr, sizeParticion, InicioParticion := ReturnMBR(pathDisco, nombreParticion)
	superBloque := SB{}
	avd := AVD{}
	dd := DD{}
	inodo := Inodo{}
	bloque := Bloque{}
	bitacora := Bitacora{}
	noEstructuras := (sizeParticion - (3 * int64(unsafe.Sizeof(superBloque)))) /
		(27 + int64(unsafe.Sizeof(avd)) + int64(unsafe.Sizeof(dd)) + (5*int64(unsafe.Sizeof(inodo)) +
			(20 * int64(unsafe.Sizeof(bloque))) + int64(unsafe.Sizeof(bitacora))))

	//NO estructuras
	var cantidadAVD int64 = noEstructuras
	var cantidadDD int64 = noEstructuras
	var cantidadInodos int64 = noEstructuras * 5
	var cantidadBloques int64 = 4 * cantidadInodos
	var Bitacoras int64 = noEstructuras
	//Bitmaps
	var InicioBitmapAVD int64 = InicioParticion + int64(unsafe.Sizeof(superBloque))
	var InicioAVD int64 = InicioBitmapAVD + cantidadAVD
	var InicioBitmapDD int64 = InicioAVD + (int64(unsafe.Sizeof(avd)) * cantidadAVD)
	var InicioDD int64 = InicioBitmapDD + cantidadDD
	var InicioBitmapInodo int64 = InicioDD + (int64(unsafe.Sizeof(dd)) * cantidadDD)
	var InicioInodo int64 = InicioBitmapInodo + cantidadInodos
	var InicioBitmapBloque int64 = InicioInodo + (int64(unsafe.Sizeof(inodo)) * cantidadInodos)
	var InicioBLoque int64 = InicioBitmapBloque + cantidadBloques
	var InicioBitacora int64 = InicioBLoque + (int64(unsafe.Sizeof(bloque)) * cantidadBloques)
	var InicioCopiaSB int64 = InicioBitacora + (int64(unsafe.Sizeof(bitacora)) * Bitacoras)

	//Inicializando SuperBloque
	copy(superBloque.Sb_nombre_hd[:], nombreDisco)
	superBloque.Sb_arbol_virtual_count = cantidadAVD
	superBloque.Sb_detalle_directorio_count = cantidadDD
	superBloque.Sb_inodos_count = cantidadInodos
	superBloque.Sb_bloques_count = cantidadBloques
	//
	superBloque.Sb_arbol_virtual_free = cantidadAVD
	superBloque.Sb_detalle_directorio_free = cantidadDD
	superBloque.Sb_inodos_free = cantidadInodos
	superBloque.Sb_bloques_free = cantidadBloques
	copy(superBloque.Sb_date_creacion[:], dt.String())
	copy(superBloque.Sb_date_ultimo_montaje[:], dt.String())
	superBloque.Sb_montajes_count = 1
	//Bitmaps
	superBloque.Sb_ap_bitmap_arbol_directorio = InicioBitmapAVD
	superBloque.Sb_ap_arbol_directorio = InicioAVD
	superBloque.Sb_ap_bitmap_detalle_directorio = InicioBitmapDD
	superBloque.Sb_ap_detalle_directorio = InicioDD
	superBloque.Sb_ap_bitmap_tabla_inodo = InicioBitmapInodo
	superBloque.Sb_ap_tabla_inodo = InicioInodo
	superBloque.Sb_ap_bitmap_bloques = InicioBitmapBloque
	superBloque.Sb_ap_bloques = InicioBLoque
	superBloque.Sb_ap_log = InicioBitacora
	superBloque.Sb_size_struct_arbol_directorio = int64(unsafe.Sizeof(avd))
	superBloque.Sb_size_struct_Detalle_directorio = int64(unsafe.Sizeof(dd))
	superBloque.Sb_size_struct_inodo = int64(unsafe.Sizeof(inodo))
	superBloque.Sb_size_struct_bloque = int64(unsafe.Sizeof(bloque))
	superBloque.Sb_first_free_bit_arbol_directorio = InicioBitmapAVD
	superBloque.Sb_first_free_bit_detalle_directoriio = InicioBitmapDD
	superBloque.Sb_dirst_free_bit_tabla_inodo = InicioBitmapInodo
	superBloque.Sb_first_free_bit_bloques = InicioBitmapBloque
	superBloque.Sb_magic_num = 61267
	superBloque.InicioCopiaSB = InicioCopiaSB
	superBloque.ConteoAVD = 0
	superBloque.ConteoDD = 0
	superBloque.ConteoInodo = 0
	superBloque.ConteoBloque = 0
	//Escribir en Particion
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()

	f.Seek(InicioParticion, 0)
	err = binary.Write(f, binary.BigEndian, &superBloque)

	f.Seek(InicioBitmapAVD, 0)
	var otro int8 = 0
	var i int64 = 0
	for i = 0; i < cantidadAVD; i++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}

	f.Seek(InicioAVD, 0)
	i = 0
	for i = 0; i < cantidadAVD; i++ {
		err = binary.Write(f, binary.BigEndian, &avd)
	}

	f.Seek(InicioBitmapDD, 0)
	i = 0
	for i = 0; i < cantidadDD; i++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}

	f.Seek(InicioDD, 0)
	i = 0
	dd.Dd_ap_detalle_directorio = -1
	for i = 0; i < cantidadDD; i++ {
		err = binary.Write(f, binary.BigEndian, &dd)
	}

	f.Seek(InicioBitmapInodo, 0)
	i = 0
	for i = 0; i < cantidadInodos; i++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}

	f.Seek(InicioInodo, 0)
	i = 0
	inodo.I_count_inodo = -1
	for i = 0; i < cantidadInodos; i++ {
		err = binary.Write(f, binary.BigEndian, &inodo)
	}

	f.Seek(InicioBitmapBloque, 0)
	i = 0
	for i = 0; i < cantidadBloques; i++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}

	f.Seek(InicioBLoque, 0)
	i = 0
	copy(bloque.Db_data[:], "")
	for i = 0; i < cantidadBloques; i++ {
		err = binary.Write(f, binary.BigEndian, &bloque)
	}

	f.Seek(InicioBitacora, 0)
	i = 0
	bitacora.Size = -1
	for i = 0; i < Bitacoras; i++ {
		err = binary.Write(f, binary.BigEndian, &bitacora)
	}

	f.Seek(InicioCopiaSB, 0)
	err = binary.Write(f, binary.BigEndian, &superBloque)

	CrearRaiz(pathDisco, InicioParticion)
	fmt.Println("SUPERBLOQUE")
	fmt.Println("Peso de SuperBloque: ", unsafe.Sizeof(superBloque))
	fmt.Println("Peso de Inodos: ", unsafe.Sizeof(inodo))
	fmt.Println("PesoBloques: ", unsafe.Sizeof(bloque))
	fmt.Println("Cantidad de Inodos: ", cantidadInodos)
	fmt.Println("CantidadBloques: ", cantidadBloques)
	fmt.Println("n: ", noEstructuras)
	fmt.Println("Particion: ", nombreParticion, NoParticion)
	fmt.Println("Peso de la particion: ", sizeParticion)
	fmt.Printf("Fecha: %s\n", mbr.MbrFechaCreacion)
	return false
}

func ReturnMBR(path string, nombreParticion string) (MBR, int64, int64) {
	mbr := MBR{}
	var Particiones [4]Particion
	var nombre2 [15]byte
	var size int64
	copy(nombre2[:], nombreParticion)
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return mbr, 0, 0
	}
	defer f.Close()

	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	Particiones = mbr.Particiones
	for i := 0; i < 4; i++ {
		if BytesNombreParticion(Particiones[i].NombreParticion) == BytesNombreParticion(nombre2) {
			size = Particiones[i].TamanioTotal
			return mbr, size, Particiones[i].Inicio_particion
		}
	}
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
			var InicioExtendida int64 = Particiones[i].Inicio_particion
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Particion_Siguiente == -1 {
				fmt.Println("No Hay particiones Logicas")
			} else {
				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Particion_Siguiente == -1 {
						break
					} else {
						f.Seek(ebr.Particion_Siguiente, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
					}
					if BytesNombreParticion(ebr.NombreParticion) == BytesNombreParticion(nombre2) {
						fmt.Println("Logica Encontrada")
						return mbr, ebr.TamanioTotal, ebr.Inicio_particion
					}

				}
			}
		}
	}
	return mbr, 0, 0
}

func CrearRaiz(pathDisco string, InicioParticion int64) bool {
	dt := time.Now()
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	f.Seek(InicioParticion, 0)
	sb := SB{}
	err = binary.Read(f, binary.BigEndian, &sb)
	/*
		Escribir 1 en bitmap avd y escribir avd
	*/
	f.Seek(sb.Sb_ap_bitmap_arbol_directorio, 0)
	var otro int8 = 0
	otro = 1
	err = binary.Write(f, binary.BigEndian, &otro)
	bitLibre, _ := f.Seek(0, os.SEEK_CUR)
	sb.Sb_first_free_bit_arbol_directorio = bitLibre
	avd := AVD{}
	copy(avd.Avd_fecha_creacion[:], dt.String())
	copy(avd.Avd_nomre_directotrio[:], "/")
	for j := 0; j < 6; j++ {
		avd.Avd_ap_array_subdirectoios[j] = -1
	}
	avd.Avd_ap_detalle_directorio = 0
	avd.Avd_ap_arbol_virtual_directorio = -1
	copy(avd.Avd_proper[:], "root")
	f.Seek(sb.Sb_ap_arbol_directorio, 0)
	err = binary.Write(f, binary.BigEndian, &avd)

	sb.Sb_arbol_virtual_free = sb.Sb_arbol_virtual_free - 1

	f.Seek(sb.Sb_ap_bitmap_detalle_directorio, 0)
	otro = 1
	err = binary.Write(f, binary.BigEndian, &otro)
	otro = 0
	bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	sb.Sb_first_free_bit_detalle_directoriio = bitLibre
	detalleDirectorio := DD{}
	arregloDD := ArregloDD{}
	copy(arregloDD.Dd_file_nombre[:], "users.txt")
	copy(arregloDD.Dd_file_date_creacion[:], dt.String())
	copy(arregloDD.Dd_file_date_modificacion[:], dt.String())
	arregloDD.Dd_file_ap_inodo = 0
	detalleDirectorio.Dd_array_files[0] = arregloDD
	detalleDirectorio.Ocupado = 1
	for j := 0; j < 5; j++ {
		if j == 0 {
			detalleDirectorio.Dd_array_files[j].Dd_file_ap_inodo = 0
		} else {
			detalleDirectorio.Dd_array_files[j].Dd_file_ap_inodo = -1
		}
	}
	detalleDirectorio.Dd_ap_detalle_directorio = -1
	f.Seek(sb.Sb_ap_detalle_directorio, 0)
	err = binary.Write(f, binary.BigEndian, &detalleDirectorio)

	sb.Sb_detalle_directorio_free = sb.Sb_detalle_directorio_free - 1
	/*
		Escribir 1 en bitmap tablaInodo y escribir Inodo
	*/
	var cantidadBloque int64 = CantidadBloqueUsar("1,G,root\n1,U,root,root,123\n")
	f.Seek(sb.Sb_ap_bitmap_tabla_inodo, 0)
	otro = 1
	err = binary.Write(f, binary.BigEndian, &otro)
	otro = 0
	bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	sb.Sb_dirst_free_bit_tabla_inodo = bitLibre
	inodo := Inodo{}
	for j := 0; j < 4; j++ {
		inodo.I_array_bloques[j] = -1
	}
	inodo.I_count_inodo = 0
	inodo.I_size_archivo = 10
	inodo.I_count_bloques_asignados = cantidadBloque
	for h := 0; h < int(cantidadBloque); h++ {
		inodo.I_array_bloques[h] = int64(h)
	}
	inodo.I_ao_indirecto = -1
	inodo.I_id_proper = 123
	f.Seek(sb.Sb_ap_tabla_inodo, 0)
	err = binary.Write(f, binary.BigEndian, &inodo)
	sb.Sb_inodos_free = sb.Sb_inodos_free - 1

	f.Seek(sb.Sb_ap_bitmap_bloques, 0)
	otro = 1
	for k := 0; k < int(cantidadBloque); k++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}
	otro = 0
	bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	sb.Sb_first_free_bit_bloques = bitLibre
	f.Seek(sb.Sb_ap_bloques, 0)
	usesTxt := []byte("1,G,root\n1,U,root,root,123\n")
	for k := 0; k < int(cantidadBloque); k++ {
		if k == 0 {
			bloque := Bloque{}
			copy(bloque.Db_data[:], string([]byte(usesTxt[0:25])))
			err = binary.Write(f, binary.BigEndian, &bloque)
		} else {
			bloque := Bloque{}
			copy(bloque.Db_data[:], string([]byte(usesTxt[k*25:len(usesTxt)])))
			err = binary.Write(f, binary.BigEndian, &bloque)
		}
		sb.Sb_bloques_free = sb.Sb_bloques_free - 1
		sb.ConteoBloque = sb.ConteoBloque + int64(k)
	}

	f.Seek(0, 0)
	f.Seek(InicioParticion, 0)
	err = binary.Write(f, binary.BigEndian, &sb)
	return false
}
func CantidadBloqueUsar(data string) int64 {
	var noBloque int64 = 0
	cont := 1
	var dataX []byte = []byte(data)
	for i := 0; i < len(dataX); i++ {
		if cont == 25 {
			noBloque = noBloque + 1
			cont = 0
		}
		cont++
	}
	if len(dataX)%25 != 0 {
		noBloque = noBloque + 1
	}
	return noBloque
}
func CantidadInodosUsar(data string) int64 {
	var noBloque int64 = 0
	cont := 0
	var dataX []byte = []byte(data)
	for i := 0; i < len(dataX); i++ {
		if cont == 25 {
			noBloque = noBloque + 1
			cont = 0
		}
		cont++
	}
	if len(dataX)%5 != 0 {
		noBloque = noBloque + 1
	}
	return noBloque
}
func DevolverSuperBlque(path string, nombreParticion string) (SB, int64) {
	mbr := MBR{}
	sb := SB{}
	var Particiones [4]Particion
	var nombre2 [15]byte
	copy(nombre2[:], nombreParticion)
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return sb, 0
	}
	defer f.Close()

	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	Particiones = mbr.Particiones
	for i := 0; i < 4; i++ {
		if BytesNombreParticion(Particiones[i].NombreParticion) == BytesNombreParticion(nombre2) {
			f.Seek(Particiones[i].Inicio_particion, 0)
			err = binary.Read(f, binary.BigEndian, &sb)
			return sb, Particiones[i].Inicio_particion
		}
	}
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
			var InicioExtendida int64 = Particiones[i].Inicio_particion
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Particion_Siguiente == -1 {
				fmt.Println("No Hay particiones Logicas")
			} else {
				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Particion_Siguiente == -1 {
						break
					} else {
						f.Seek(ebr.Particion_Siguiente, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
					}
					if BytesNombreParticion(ebr.NombreParticion) == BytesNombreParticion(nombre2) {
						fmt.Println("Logica Encontrada")
						f.Seek(ebr.Inicio_particion, 0)
						err = binary.Read(f, binary.BigEndian, &sb)
						return sb, ebr.Inicio_particion
					}

				}
			}
		}
	}
	return sb, 0
}

func EjecutarComandoLogin(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (bool, string) {
	ParamValidos := true
	usuario := ""
	var propiedades [3]string

	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]

			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-usuario":
				propiedades[0] = propiedadTemp.Val
			case "-password":
				propiedades[1] = string(propiedadTemp.Val)
			case "-id":
				propiedades[2] = propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		fmt.Println("usuario", propiedades[0])
		fmt.Println("password", propiedades[1])
		fmt.Println("id", propiedades[2])

		ParamValidos, usuario = ExecuteLogin(propiedades[0], propiedades[1], propiedades[2], ListaDiscos)
		return ParamValidos, usuario
	} else {
		ParamValidos = false
		return ParamValidos, usuario
	}
}
func ExecuteLogin(usuario string, password string, id string, ListaDiscos *list.List) (bool, string) {
	idValido := IdValido(id, ListaDiscos)
	fmt.Println(idValido)
	if idValido == false {
		fmt.Println("El id ingresado no esta montada en el disco.")
		return false, ""
	} else if global != "" {
		fmt.Println("Ya estas en una sesión ")
		return false, ""
	}
	pathDisco, nombreParticion, nombreDisco := RecorrerListaDisco(id, ListaDiscos)
	mbr, sizeParticion, InicioParticion := ReturnMBR(pathDisco, nombreParticion)
	superBloque := SB{}
	f, err := os.OpenFile(pathDisco, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta " + pathDisco)
		return false, ""
	}
	defer f.Close()
	f.Seek(InicioParticion, 0)
	err = binary.Read(f, binary.BigEndian, &superBloque)
	/*
	   Obtener avd raiz
	*/
	avd := AVD{}
	dd := DD{}
	inodo := Inodo{}
	bloque := Bloque{}
	f.Seek(superBloque.Sb_ap_arbol_directorio, 0)
	err = binary.Read(f, binary.BigEndian, &avd)
	apuntadorDD := avd.Avd_ap_detalle_directorio
	f.Seek(superBloque.Sb_ap_detalle_directorio, 0)
	for i := 0; i < int(superBloque.Sb_arbol_virtual_free); i++ {
		err = binary.Read(f, binary.BigEndian, &dd)
		if i == int(apuntadorDD) {
			break
		}
	}
	apuntadorInodo := dd.Dd_array_files[0].Dd_file_ap_inodo
	f.Seek(superBloque.Sb_ap_tabla_inodo, 0)
	for i := 0; i < int(superBloque.Sb_inodos_free); i++ {
		err = binary.Read(f, binary.BigEndian, &inodo)
		if i == int(apuntadorInodo) {
			break
		}
	}
	var userstxt string = ""

	//Leer Users.txt
	posicion := 0
	f.Seek(superBloque.Sb_ap_bloques, 0)
	for i := 0; i < int(superBloque.Sb_inodos_free); i++ {
		err = binary.Read(f, binary.BigEndian, &bloque)

		if int(inodo.I_array_bloques[posicion]) != -1 && int(inodo.I_array_bloques[posicion]) == i {
			userstxt += ConvertData(bloque.Db_data)
		} else if int(inodo.I_array_bloques[posicion]) == -1 {
			break
		} else {
			break
		}
		if posicion < 4 {
			posicion++
		} else if posicion == 4 {
			posicion = 0
		}
	}
	lineaUsuarioTxt := strings.Split(userstxt, "\n")
	for i := 0; i < len(lineaUsuarioTxt); i++ {
		if len(lineaUsuarioTxt[i]) != 17 {
			usuario_grupo := strings.Split(lineaUsuarioTxt[i], ",")
			if usuario_grupo[1] == "U" {
				if usuario_grupo[3] == usuario && usuario_grupo[4] == password {
					fmt.Println("Inicio sesión correctamente")
					return true, usuario
				}
			}
		}
	}
	fmt.Println(nombreDisco, mbr.MbrTamanio, sizeParticion)
	return false, ""
}

func RecorrerListaDisco(id string, ListaDiscos *list.List) (string, string, string) {
	Id := strings.ReplaceAll(id, "14", "")
	//NoParticion := Id[1:]
	IdDisco := Id[:1]
	pathDisco := ""
	nombreParticion := ""
	nombreDisco := ""
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if BytesToString(disco.Id) == IdDisco {
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if mountTemp.Id == id {
					copy(mountTemp.EstadoMKS[:], "1")
					nombreParticion = mountTemp.NombreParticion
					pathDisco = disco.Path
					nombreDisco = disco.NombreDisco
					return pathDisco, nombreParticion, nombreDisco
					break
				}
			}

		}
		element.Value = disco
	}
	return "", "", ""
}

type Propiedad struct {
	Name string
	Val  string
}
type Comando struct {
	Name        string
	Propiedades []Propiedad
}

type Particion struct {
	Status_particion [1]byte
	TipoParticion    [1]byte
	TipoAjuste       [2]byte
	Inicio_particion int64
	TamanioTotal     int64
	NombreParticion  [15]byte
}

type MBR struct {
	MbrTamanio       int64
	MbrFechaCreacion [19]byte
	NoIdentificador  int64
	Particiones      [4]Particion
}

type EBR struct {
	Status_particion    [1]byte
	TipoAjuste          [2]byte
	Inicio_particion    int64
	Particion_Siguiente int64
	TamanioTotal        int64
	NombreParticion     [15]byte
}

type MOUNT struct {
	NombreParticion string
	Id              string
	Estado          [1]byte
	EstadoMKS       [1]byte
}

type DISCO struct {
	NombreDisco string
	Path        string
	Id          [1]byte
	Estado      [1]byte
	Particiones [100]MOUNT
}

type Integers struct {
	I1  uint16
	I2  int32
	I3  int64
	DOS byte
}

type SB struct {
	Sb_nombre_hd                          [15]byte
	Sb_arbol_virtual_count                int64
	Sb_detalle_directorio_count           int64
	Sb_inodos_count                       int64
	Sb_bloques_count                      int64
	Sb_arbol_virtual_free                 int64
	Sb_detalle_directorio_free            int64
	Sb_inodos_free                        int64
	Sb_bloques_free                       int64
	Sb_date_creacion                      [19]byte
	Sb_date_ultimo_montaje                [19]byte
	Sb_montajes_count                     int64
	Sb_ap_bitmap_arbol_directorio         int64
	Sb_ap_arbol_directorio                int64
	Sb_ap_bitmap_detalle_directorio       int64
	Sb_ap_detalle_directorio              int64
	Sb_ap_bitmap_tabla_inodo              int64
	Sb_ap_tabla_inodo                     int64
	Sb_ap_bitmap_bloques                  int64
	Sb_ap_bloques                         int64
	Sb_ap_log                             int64
	Sb_size_struct_arbol_directorio       int64
	Sb_size_struct_Detalle_directorio     int64
	Sb_size_struct_inodo                  int64
	Sb_size_struct_bloque                 int64
	Sb_first_free_bit_arbol_directorio    int64
	Sb_first_free_bit_detalle_directoriio int64
	Sb_dirst_free_bit_tabla_inodo         int64
	Sb_first_free_bit_bloques             int64
	Sb_magic_num                          int64
	InicioCopiaSB                         int64
	ConteoAVD                             int64
	ConteoDD                              int64
	ConteoInodo                           int64
	ConteoBloque                          int64
}

type ArregloDD struct {
	Dd_file_nombre            [15]byte
	Dd_file_ap_inodo          int64
	Dd_file_date_creacion     [19]byte
	Dd_file_date_modificacion [19]byte
}
type DD struct {
	Dd_array_files           [5]ArregloDD
	Dd_ap_detalle_directorio int64
	Ocupado                  int8
}

type Inodo struct {
	I_count_inodo             int64
	I_size_archivo            int64
	I_count_bloques_asignados int64
	I_array_bloques           [4]int64
	I_ao_indirecto            int64
	I_id_proper               int64
}

type AVD struct {
	Avd_fecha_creacion              [19]byte
	Avd_nomre_directotrio           [15]byte
	Avd_ap_array_subdirectoios      [6]int64
	Avd_ap_detalle_directorio       int64
	Avd_ap_arbol_virtual_directorio int64
	Avd_proper                      [10]byte
}

type Bloque struct {
	Db_data [25]byte
}

type Bitacora struct {
	Log_tipo_operacion [19]byte
	Log_tipo           [1]byte
	Log_nombre         [35]byte
	Log_Contenido      [25]byte
	Log_fecha          [19]byte
	Size               int64
}

func BytesNombreParticion(data [15]byte) string {
	return string(data[:])
}
func ConvertData(data [25]byte) string {
	return string(data[:])
}
